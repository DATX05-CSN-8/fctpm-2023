package tpmpool

import (
	"errors"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/tpminstantiator"
)

type tpmallocator interface {
	Allocate() (*tpminstantiator.TpmInstance, error)
	Return(*tpminstantiator.TpmInstance) error
}

type tpmPoolService struct {
	alloc  tpmallocator
	readyq []*tpminstantiator.TpmInstance
	mu     chan int
}

func NewTpmPoolService(basepath string, size int) (*tpmPoolService, error) {
	t := tpmPoolService{
		readyq: make([]*tpminstantiator.TpmInstance, 0),
		alloc:  tpminstantiator.NewTpmInstantiatorServiceWithBasePath(basepath),
		mu:     make(chan int, 1),
	}
	t.mu <- 1
	defer func() { <-t.mu }()
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
	s.mu <- 1
	defer func() { <-s.mu }()
	if len(s.readyq) < 0 {
		return nil, errors.New("Error occurred, no more elements in tpm pool ready queue")
	}
	tpminstance := s.readyq[0]
	s.readyq = s.readyq[1:]
	return tpminstance, nil
}

func (s *tpmPoolService) Return(instance *tpminstantiator.TpmInstance) error {
	err := s.alloc.Return(instance)
	if err != nil {
		return err
	}
	return nil
}
