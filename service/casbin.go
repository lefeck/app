package service

import (
	"app/forms"
	"app/model"
	"app/repository"
)

type CasbinService struct {
	casbinRepository repository.ICasbinRepository
}

type ICasbinService interface {
	//Create(param *forms.CasbinCreateRequest) error
	Create(casbin *model.Casbin) (*model.Casbin, error)
	List(param *forms.CasbinListRequest) [][]string
}

func NewCasbinService(casbinRepository repository.ICasbinRepository) ICasbinService {
	return &CasbinService{casbinRepository}
}

//func (s *CasbinService) Create(param *forms.CasbinCreateRequest) error {
//
//	for _, v := range param.CasbinInfos {
//		err := s.casbinRepository.CasbinCreate(param.RoleId, v.Path, v.Method)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (c *CasbinService) Create(casbin *model.Casbin) (*model.Casbin, error) {
	//cm := model.Casbin{
	//	PType:  "p",
	//	RoleId: roleId,
	//	Path:   path,
	//	Method: method,
	//}
	return c.casbinRepository.Create(casbin)
}

func (s *CasbinService) List(param *forms.CasbinListRequest) [][]string {
	return s.casbinRepository.CasbinList(param.RoleID)
}
