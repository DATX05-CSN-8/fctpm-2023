package resourcepool

import (
	"errors"
	"sync"
)

type resourceAllocator[T any] interface {
	Allocate() (*T, error)
	Return(*T) error
}

type resourcePool[T any] struct {
	alloc  resourceAllocator[T]
	readyq []*T
	mu     sync.Mutex
}

func NewResourcePool[T any](size int, allocator resourceAllocator[T]) (*resourcePool[T], error) {
	t := resourcePool[T]{
		readyq: make([]*T, size),
		alloc:  allocator,
		mu:     sync.Mutex{},
	}
	for i := 0; i < size; i++ {
		inst, err := t.alloc.Allocate()
		if err != nil {
			return nil, err
		}
		t.readyq[i] = inst
	}
	return &t, nil
}

func (p *resourcePool[T]) Allocate() (*T, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.readyq) <= 0 {
		return nil, errors.New("Error occurred, no more elements in tpm pool ready queue")
	}
	inst := p.readyq[0]
	p.readyq = p.readyq[1:]
	return inst, nil
}

func (p *resourcePool[T]) Return(instance *T) error {
	err := p.alloc.Return(instance)
	return err
}
