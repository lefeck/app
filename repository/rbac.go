package repository

import (
	"app/database"
	"app/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type rbacRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

func (r *rbacRepository) List() ([]model.Role, error) {
	roles := []model.Role{}
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *rbacRepository) ListResources() ([]model.Resource, error) {
	resources := make([]model.Resource, 0)
	if err := r.db.Order("name").Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *rbacRepository) Create(role *model.Role) (*model.Role, error) {
	if err := r.db.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *rbacRepository) CreateResource(resource *model.Resource) (*model.Resource, error) {
	if err := r.db.Create(resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func (r *rbacRepository) CreateResources(resources []model.Resource, conds ...clause.Expression) error {
	//r.db.Clauses(clause.OnConflict{})
	return r.db.Clauses(conds...).Create(resources).Error
}

func (r *rbacRepository) GetRoleByID(id int) (*model.Role, error) {
	role := &model.Role{}
	if err := r.db.First(role, id).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *rbacRepository) GetResource(id int) (*model.Resource, error) {
	resoure := &model.Resource{}
	if err := r.db.First(resoure, id).Error; err != nil {
		return nil, err
	}
	return resoure, nil
}

func (r *rbacRepository) GetRoleByName(name string) (*model.Role, error) {
	role := &model.Role{}
	if err := r.db.Preload(model.UserAssociation).Where("name = ?", name).First(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *rbacRepository) Update(role *model.Role) (*model.Role, error) {
	if err := r.db.Updates(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *rbacRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Role{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *rbacRepository) DeleteResource(id uint) error {
	return r.db.Delete(&model.Resource{}, id).Error
}

func (r *rbacRepository) Migrate() error {
	return r.db.AutoMigrate(&model.Role{}, &model.Resource{})
}

func newRbacRespository(db *gorm.DB, rdb *database.RedisDB) RBACRepository {
	return &rbacRepository{db: db, rdb: rdb}
}
