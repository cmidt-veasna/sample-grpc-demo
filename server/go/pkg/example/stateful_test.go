package example

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
)

var numEle = 4

func TestLoadingPersistedData(t *testing.T) {
	filename = "test.bin"
	os.Remove(filename)
	ch := make(chan *Element, numEle)
	cl := make(chan bool, 1)
	go statefulDatabaseClosable(ch, cl)
	expectedSize := int64(0)
	for i := 1; i <= numEle; i++ {
		ele := &Element{
			Id:        fmt.Sprintf("id-%0d", i),
			Name:      fmt.Sprintf("sample-%0d", i),
			Age:       int32(10 + i),
			CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
			UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		}
		ch <- ele
		b, _ := proto.Marshal(ele)
		expectedSize += int64(len(b) + 4)
	}
	close(ch)
	<-cl
	ff, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		t.Error("unable to locate database file", filename, "due to", err)
	} else if ff == nil {
		t.Error("unable to locate database file", filename, "state file nil")
	} else if size := ff.Size(); size != expectedSize {
		t.Error("wrong database file", filename, "size expect", expectedSize, "but got", size)
	}
}

func TestStatefulDatabase(t *testing.T) {
	filename = "test.bin"
	ff, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		t.Error("unable to locate database file", filename, "due to", err)
		return
	} else if ff == nil {
		t.Error("unable to locate database file", filename, "state file nil")
	}

	eles := loadPersistData()
	if size := len(eles); size != numEle {
		t.Error("unexpected total number of element in database file", filename, "expected", numEle, "but got", size)
	}
	for k, v := range eles {
		if v.Id == "" || v.Name == "" || v.Age == 0 || v.CreatedAt == "" || v.UpdatedAt == "" {
			t.Error("element properties is empty")
		}
		if k != v.Id {
			t.Error("wrong id, expected", k, "but got", v.Id)
		}
	}
	os.Remove(filename)
}

