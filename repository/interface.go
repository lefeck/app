package repository

import (
	"app/model"
	"context"
	"gorm.io/gorm/clause"
)

// 工厂模式接口
type Repository interface {
	User() UserRepository
	Article() ArticleRepository
	Category() CategoryRepository
	Comment() CommentRepository
	Tag() TagRepository
	Like() LikeRepository
	Close() error
	Ping(ctx context.Context) error
	//Init() error
	Migrant
	//Group() GroupRepository
	//RBAC() RBACRepository
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

// Article 接口
type ArticleRepository interface {
	GetArticleByID(uint) (*model.Article, error)
	GetArticleByName(string) (*model.Article, error)
	List() ([]model.Article, error)
	Create(*model.User, *model.Article) (*model.Article, error)
	Update(*model.Article) (*model.Article, error)
	Delete(uint) error
	IncView(id uint) error
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

type CategoryRepository interface {
	GetCategoriesByArticle(*model.Article) ([]model.Category, error)
	Delete(cid uint) error
	Create(category *model.Category) (*model.Category, error)
	Update(category *model.Category) (*model.Category, error)
}

type TagRepository interface {
	GetTagsByArticle(article *model.Article) ([]model.Tag, error)
	Create(tag *model.Tag) (*model.Tag, error)
	Delete(id uint) error
	List() ([]model.Tag, error)
	Update(tag *model.Tag) (*model.Tag, error)
}

type CommentRepository interface {
	Add(comment *model.Comment) (*model.Comment, error)
	Delete(id string) error
	List(aid string) ([]model.Comment, error)
}

type LikeRepository interface {
	Add(aid, uid uint) error
	Delete(aid, uid uint) error
	Get(aid, uid uint) (bool, error)
	GetLikeByUser(uid uint) ([]model.Like, error)
}
