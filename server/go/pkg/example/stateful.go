package example

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
)

var filename = "/sample.data.bin"

// load data from file
func loadPersistData() (m map[string]*Element) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return make(map[string]*Element)
	}
	m = make(map[string]*Element)
	r := bytes.NewReader(b)
	size := int32(0)
	start := int32(0)
	totalBuf := len(b)
	for {
		if err = binary.Read(r, binary.LittleEndian, &size); err != nil {
			if err == io.EOF {
				break
			}
			log.Println("warning unable to decode underlying data from stateful database file!", err)
			return
		}
		if int(start) >= totalBuf || int(start+size) > totalBuf {
			log.Println("warning malformed underlying data from stateful database file!", err)
			return
		}
		ele := &Element{}
		// increase 4 bytes for data size
		start += 4
		if err = proto.Unmarshal(b[start:start+size], ele); err != nil {
			log.Println("warning unable to decode protobuf data from stateful database file!")
			return
		}
		start += size
		r.Seek(int64(size), io.SeekCurrent)
		m[ele.Id] = ele
	}
	return
}

func statefulDatabase(ch <-chan *Element) {
	statefulDatabaseClosable(ch, nil)
}

// loop queue save file
func statefulDatabaseClosable(ch <-chan *Element, cl chan<- bool) {
	var f *os.File
	for {
		select {
		case ele, ok := <-ch:
			if !ok {
				if f != nil {
					f.Close()
				}
				// channel closed
				if cl != nil {
					cl <- true
				}
				return
			}
			if b, err := proto.Marshal(ele); err != nil {
				log.Println("warning unable to encode elements", err)
			} else {
				if f == nil {
					retryCount := 0
				Retry:
					if f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0700); err != nil {
						// retry 5 time
						if retryCount < 5 {
							log.Println("warning unable to open database file", err)
							retryCount++
							goto Retry
						} else {
							log.Fatal("warning unable to open database file", err)
						}
					}
				}
				if err = binary.Write(f, binary.LittleEndian, int32(len(b))); err != nil {
					log.Fatal("warning unable to persist elements meta", err)
					return
				}
				if _, err = f.Write(b); err != nil {
					log.Fatal("warning unable to persist elements", err)
				} else {
					log.Println("elements saved on disk")
				}
			}
		}
	}
}
