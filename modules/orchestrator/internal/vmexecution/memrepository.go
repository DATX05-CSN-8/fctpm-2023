package vmexecution

import (
	"fmt"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type maprepository struct {
	backend cmap.ConcurrentMap[string, *VMExecution]
}

func NewMapRepository() *maprepository {
	return &maprepository{
		backend: cmap.New[*VMExecution](),
	}
}

func (r *maprepository) FindAll() ([]VMExecution, error) {
	var all []VMExecution
	for item := range r.backend.IterBuffered() {
		all = append(all, *item.Val)
	}
	return all, nil
}

func (r *maprepository) Create(e VMExecution) (VMExecution, error) {
	r.backend.Set(e.Id, &e)
	return e, nil
}

func (r *maprepository) FindById(id string) (VMExecution, error) {
	e, ok := r.backend.Get(id)
	if !ok {
		return *e, fmt.Errorf("Key %s not found", id)
	}
	return *e, nil
}

func (r *maprepository) Update(e *VMExecution) (*VMExecution, error) {
	vmi, err := r.Create(*e)
	return &vmi, err
}

func (r *maprepository) Delete(e VMExecution) error {
	r.backend.Remove(e.Id)
	return nil
}
