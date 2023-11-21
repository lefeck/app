package service

import (
	"app/model"
	"app/repository"
	"app/utils/request"
	"gorm.io/gorm/clause"
	"strconv"
)

type rbacService struct {
	rbacRepository repository.RBACRepository
}

func NewrbacService(rbacRepository repository.RBACRepository) RBACService {
	return &rbacService{
		rbacRepository: rbacRepository,
	}
}

func (r *rbacService) List() ([]model.Role, error) {
	return r.rbacRepository.List()
}

func (r *rbacService) ListResources() ([]model.Resource, error) {
	return r.rbacRepository.ListResources()
}

func (r *rbacService) Create(role *model.Role) (*model.Role, error) {
	return r.rbacRepository.Create(role)
}

func (r *rbacService) CreateResource(resource *model.Resource) (*model.Resource, error) {
	return r.rbacRepository.CreateResource(resource)
}

func (r *rbacService) CreateResources(resources []model.Resource, conds ...clause.Expression) error {
	return r.rbacRepository.CreateResources(resources, conds...)
}

func (r *rbacService) GetRoleByID(id string) (*model.Role, error) {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return r.rbacRepository.GetRoleByID(rid)
}

func (r *rbacService) GetResource(id string) (*model.Resource, error) {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return r.rbacRepository.GetResource(rid)
}

func (r *rbacService) GetRoleByName(name string) (*model.Role, error) {
	return r.GetRoleByName(name)
}

func (r *rbacService) Update(id string, role *model.Role) (*model.Role, error) {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	role.ID = uint(rid)
	return r.rbacRepository.Update(role)
}

func (r *rbacService) Delete(id string) error {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return r.rbacRepository.Delete(uint(rid))
}

func (r *rbacService) DeleteResource(id string) error {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return r.rbacRepository.DeleteResource(uint(rid))
}

func (r *rbacService) ListOperations() ([]model.Operation, error) {
	return []model.Operation{
		model.AllOperation,
		model.EditOperation,
		model.ViewOperation,
		request.CreateOperation,
		request.PatchOperation,
		request.UpdateOperation,
		request.GetOperation,
		request.ListOperation,
		request.DeleteOperation,
		"log",
		"exec",
		"proxy",
	}, nil
}
