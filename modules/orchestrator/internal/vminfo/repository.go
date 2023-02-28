package vminfo

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]VMInfo, error)
	FindById(id string) (VMInfo, error)
	Create(data VMInfo) (VMInfo, error)
	Update(data VMInfo) (VMInfo, error)
	Delete(data VMInfo) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]VMInfo, error) {
	var all []VMInfo
	err := r.db.Find(&all).Error
	return all, err
}

func (r *repository) Create(vmInfo VMInfo) (VMInfo, error) {
	err := r.db.Create(&vmInfo).Error
	return vmInfo, err
}

func (r *repository) FindById(id string) (VMInfo, error) {
	var vmInfo VMInfo
	err := r.db.Where("id = ?", id).First(&vmInfo).Error
	return vmInfo, err
}

func (r *repository) Update(vmInfo VMInfo) (VMInfo, error) {
	err := r.db.Save(&vmInfo).Error
	return vmInfo, err
}

func (r *repository) Delete(vmInfo VMInfo) error {
	err := r.db.Delete(&vmInfo).Error
	return err
}
