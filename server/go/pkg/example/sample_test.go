package example

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
)

var elems []*Element

// implement equal method
func (e *Element) equal(ele *Element) bool {
	// skip updatedAt as it may have change
	return e.Id == ele.Id && e.Status == ele.Status &&
		e.Age == ele.Age && e.CreatedAt == ele.CreatedAt
}

// implement available in method
func (e *Element) availableIn(eles []*Element) bool {
	for i := range eles {
		if e.equal(eles[i]) {
			return true
		}
	}
	return false
}

func TestPersistElement(t *testing.T) {
	filename = "test.bin"
	s := New()
	ctx := context.Background()
	for i := 0; i < numEle; i++ {
		ele, err := s.PersistElement(ctx, &Element{
			Name:   fmt.Sprintf("test %d", i+1),
			Age:    int32(42 + i),
			Status: uint32(i + 1),
		})
		if err != nil {
			t.Error("test index", i, "unexpected error when persist element", err)
		} else if ele.Id == "" {
			t.Error("test index", i, "expected element id to be available but given empty")
		}
		elems = append(elems, ele)
	}
	close(s.ch)
	// ensure that channel is closed and file is closed
	<-s.cl
	// reset singleton
	sample = nil
	once = sync.Once{}
}

func TestListElement(t *testing.T) {
	filename = "test.bin"
	s := New()
	ctx := context.Background()

	efs := []*ElementFilter{
		{Id: elems[0].Id}, {Id: "-----"}, {Name: "test%"}, {Age: "[30, 44]"}, {Status: "{1,3,4}"},
	}
	expected := []int{1, 0, 4, 3, 3}
	equals := [][]*Element{
		elems[:1], {}, elems, elems[:3], append(elems[:1], elems[1:]...),
	}

	for i, elef := range efs {
		eles, err := s.ListElement(ctx, elef)
		if err != nil {
			t.Error("test index", i, "unexpected error when list elements", err)
		}
		if eles == nil {
			t.Error("test index", i, "unexpected function return elements, it's nil")
			continue
		}
		if size := len(eles.Elements); size != expected[i] {
			t.Error("test index", i, "expected number element returned", expected[i], "but got", size)
		}
		// test object is equal
		for in := range eles.Elements {
			if !eles.Elements[in].availableIn(equals[i]) {
				t.Error("test index", i, "returned element expected in", equals[i], "but got", eles.Elements[in])
			}
		}
	}
	close(s.ch)
	// ensure that channel is closed and file is closed
	<-s.cl
	// remove test database file
	os.Remove(filename)
}
