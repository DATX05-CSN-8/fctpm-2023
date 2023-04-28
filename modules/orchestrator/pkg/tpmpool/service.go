package tpmpool

import (
	"errors"
	"sync"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/tpminstantiator"
)

type tpmallocator interface {
	Allocate() (*tpminstantiator.TpmInstance, error)
	Return(*tpminstantiator.TpmInstance) error
}

type tpmPoolService struct {
	alloc  tpmallocator
	readyq []*tpminstantiator.TpmInstance
	//busyq  []*tpminstantiator.TpmInstance
	mu sync.Mutex
}

func NewTpmPoolService(basepath string, size int) (*tpmPoolService, error) {
	t := tpmPoolService{
		readyq: make([]*tpminstantiator.TpmInstance, 0),
		alloc:  tpminstantiator.NewTpmInstantiatorServiceWithBasePath(basepath),
		//busyq:  make([]*tpminstantiator.TpmInstance, 0),
		mu: sync.Mutex{},
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	for i := 0; i < size; i++ {
		tpminstance, err := t.alloc.Allocate()
		if err != nil {
			return nil, err
		}
		t.readyq = append(t.readyq, tpminstance)
	}
	return &t, nil
}

func (s *tpmPoolService) Allocate() (*tpminstantiator.TpmInstance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.readyq) < 0 {
		return nil, errors.New("Error occurred, no more elements in readyq")
	}
	tpminstance := s.readyq[0]
	s.readyq = s.readyq[1:]
	//s.busyq = append(s.busyq, tpminstance)
	return tpminstance, nil
}

func (s *tpmPoolService) Return(instance *tpminstantiator.TpmInstance) error {
	//s.busyq = pop_spec_(s.busyq, tpminstance)
	err := s.alloc.Return(instance)
	if err != nil {
		return err
	}
	return nil
}
