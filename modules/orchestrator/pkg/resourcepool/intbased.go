package resourcepool

import (
	"errors"
	"sync/atomic"
)

type altResourcePool[T any] struct {
	alloc  resourceAllocator[T]
	readyq []*T
	curr   int64
}

func NewAltResourcePool[T any](size int, allocator resourceAllocator[T]) (*altResourcePool[T], error) {
	t := altResourcePool[T]{
		readyq: make([]*T, size),
		alloc:  allocator,
		curr:   -1,
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

func (p *altResourcePool[T]) Allocate() (*T, error) {
	i := atomic.AddInt64(&p.curr, 1)

	if int(i) >= len(p.readyq) {
		return nil, errors.New("Error occurred, no more elements in tpm pool ready queue")
	}
	return p.readyq[i], nil
}

func (p *altResourcePool[T]) Return(instance *T) error {
	err := p.alloc.Return(instance)
	return err
}
