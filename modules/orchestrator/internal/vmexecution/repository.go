package vmexecution

import "fmt"

type Repository interface {
	FindAll() ([]VMExecution, error)
	FindById(id string) (VMExecution, error)
	Create(data VMExecution) (VMExecution, error)
	Update(data VMExecution) (VMExecution, error)
	Delete(data VMExecution) error
}

type repository struct {
	data map[string]VMExecution
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]VMExecution),
	}
}

func (r *repository) FindAll() ([]VMExecution, error) {
	var ret []VMExecution

	for _, e := range r.data {
		ret = append(ret, e)
	}
	return ret, nil
}

func (r *repository) Create(e VMExecution) (VMExecution, error) {
	r.data[e.Id] = e
	return e, nil
}

func (r *repository) FindById(id string) (VMExecution, error) {
	d, prs := r.data[id]
	if !prs {
		return d, fmt.Errorf("NOT_FOUND")
	}
	return r.data[id], nil
}

func (r *repository) Update(e VMExecution) (VMExecution, error) {
	return r.Create(e)
}

func (r *repository) Delete(e VMExecution) error {
	delete(r.data, e.Id)
	return nil
}
