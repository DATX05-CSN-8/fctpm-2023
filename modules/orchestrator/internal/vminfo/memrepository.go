package vminfo

import (
	"fmt"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type maprepository struct {
	backend cmap.ConcurrentMap[string, *VMInfo]
}

func NewMapRepository() *maprepository {
	return &maprepository{
		backend: cmap.New[*VMInfo](),
	}
}

func (r *maprepository) FindAll() ([]VMInfo, error) {
	var all []VMInfo
	for item := range r.backend.IterBuffered() {
		all = append(all, *item.Val)
	}
	return all, nil
}

func (r *maprepository) Create(vmInfo VMInfo) (VMInfo, error) {
	r.backend.Set(vmInfo.Id, &vmInfo)
	return vmInfo, nil
}

func (r *maprepository) FindById(id string) (VMInfo, error) {
	vminfo, ok := r.backend.Get(id)
	if !ok {
		return *vminfo, fmt.Errorf("Key %s not found", id)
	}
	return *vminfo, nil
}

func (r *maprepository) Update(vmInfo *VMInfo) (*VMInfo, error) {
	vmi, err := r.Create(*vmInfo)
	return &vmi, err
}

func (r *maprepository) Delete(vmInfo VMInfo) error {
	r.backend.Remove(vmInfo.Id)
	return nil
}
