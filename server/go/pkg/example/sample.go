package example

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"example.com/envoygrpc/pkg/query"
)

var sample *Sample
var once sync.Once
var lock sync.RWMutex

type Sample struct {
	m       map[string]*Element
	ch      chan *Element
	cl      chan bool
	counter uint32
}

func New() *Sample {
	// load data from file
	once.Do(func() {
		sample = &Sample{
			m:  loadPersistData(),
			ch: make(chan *Element, 10), // 10 element queue
			cl: make(chan bool, 1),
		}
		go statefulDatabaseClosable(sample.ch, sample.cl)
	})
	return sample
}

func (s *Sample) newId() string {
	lock.Lock()
	defer lock.Unlock()
	return fmt.Sprintf("ele-%d-%d", time.Now().UTC().UnixNano(), atomic.AddUint32(&s.counter, 1))
}

func (s *Sample) PersistElement(ctx context.Context, ele *Element) (result *Element, err error) {
	ele.Id = s.newId()
	ele.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	ele.CreatedAt = ele.UpdatedAt
	s.m[ele.Id] = ele
	result = ele
	// send element to channel
	s.ch <- ele
	return
}

func (s *Sample) ListElement(ctx context.Context, eleFilter *ElementFilter) (results *Elements, err error) {
	actions := make([]query.Action, 5)
	actions[0] = query.FilterActionString(eleFilter.Name)
	if actions[1], err = query.FilterActionNumber(eleFilter.Age); err != nil {
		err = status.Error(codes.InvalidArgument, fmt.Sprintf("filter value age invalid %s", err.Error()))
		return
	}
	if actions[2], err = query.FilterActionNumber(eleFilter.Status); err != nil {
		err = status.Error(codes.InvalidArgument, fmt.Sprintf("filter value status invalid %s", err.Error()))
		return
	}
	if actions[3], err = query.FilterActionDateTime(eleFilter.CreatedAt); err != nil {
		err = status.Error(codes.InvalidArgument, fmt.Sprintf("filter value create date invalid %s", err.Error()))
		return
	}
	if actions[4], err = query.FilterActionDateTime(eleFilter.UpdatedAt); err != nil {
		err = status.Error(codes.InvalidArgument, fmt.Sprintf("filter value update date invalid %s", err.Error()))
		return
	}

	cond := false
	results = &Elements{Elements: make([]*Element, 0)}
	for k, v := range s.m {
		if eleFilter.Id != "" && k != eleFilter.Id {
			continue
		}
		cond = true
		cond = cond && (actions[0] == nil || actions[0](v.Name))
		cond = cond && (actions[1] == nil || actions[1](int64(v.Age)))
		cond = cond && (actions[2] == nil || actions[2](int64(v.Status)))
		cond = cond && (actions[3] == nil || actions[3](v.CreatedAt))
		cond = cond && (actions[4] == nil || actions[4](v.UpdatedAt))
		if cond {
			results.Elements = append(results.Elements, v)
		}
	}
	return
}
