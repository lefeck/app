package repository

import (
	"app/model"
	"context"
	"gorm.io/gorm/clause"
)

// 工厂模式接口
type Repository interface {
	User() UserRepository
	//Group() GroupRepository
	Post() PostRepository
	//RBAC() RBACRepository
	Close() error
	Ping(ctx context.Context) error
	//Init() error
	Migrant
}

type Migrant interface {
	Migrate() error
}

// user实现的接口
type UserRepository interface {
	GetUserByID(uint) (*model.User, error)
	GetUserByAuthID(authType, authID string) (*model.User, error)
	GetUserByName(string) (*model.User, error)
	//List() ([]model.User, error)
	List(pageSize int, pageNum int) (int, []interface{})
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	FindAll(userlist []model.User) ([]model.User, error)
	Delete(id int) error
	Migrate() error
}

// post 接口
type PostRepository interface {
	GetPostByID(uint) (*model.Post, error)
	GetPostByName(string) (*model.Post, error)
	List() ([]model.Post, error)
	Create(*model.User, *model.Post) (*model.Post, error)
	Update(*model.Post) (*model.Post, error)
	Delete(uint) error
	GetTags(*model.Post) ([]model.Tag, error)
	GetCategories(*model.Post) ([]model.Category, error)
	IncView(id uint) error
	AddLike(pid, uid uint) error
	DelLike(pid, uid uint) error
	GetLike(pid, uid uint) (bool, error)
	GetLikeByUser(uid uint) ([]model.Like, error)
	AddComment(comment *model.Comment) (*model.Comment, error)
	DelComment(id string) error
	ListComment(pid string) ([]model.Comment, error)
	Migrate() error
}

type RBACRepository interface {
	List() ([]model.Role, error)
	ListResources() ([]model.Resource, error)
	Create(role *model.Role) (*model.Role, error)
	CreateResource(resource *model.Resource) (*model.Resource, error)
	CreateResources(resources []model.Resource, conds ...clause.Expression) error
	GetRoleByID(id int) (*model.Role, error)
	GetResource(id int) (*model.Resource, error)
	GetRoleByName(name string) (*model.Role, error)
	Update(role *model.Role) (*model.Role, error)
	Delete(id uint) error
	DeleteResource(id uint) error
	Migrate() error
}
