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
	Like() likeRepository
	Tag() TagRepository
	Close() error
	Ping(ctx context.Context) error
	//Init() error
	//Group() GroupRepository
	//RBAC() RBACRepository
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

// article 接口
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

type CategoryRepository interface {
	GetCategories(article *model.Article) ([]model.Category, error)
	Create(*model.Category) (*model.Category, error)
	Delete(cid uint) error
	Update(*model.Category) (*model.Category, error)
	List() ([]model.Category, error)
}

type LikeRepository interface {
	AddLike(aid, uid uint) error
	DelLike(aid, uid uint) error
	GetLike(aid, uid uint) (bool, error)
	GetLikeByUser(uid uint) ([]model.Like, error)
}

type CommentRepository interface {
	AddComment(comment *model.Comment) (*model.Comment, error)
	DelComment(id string) error
	ListComment(aid string) ([]model.Comment, error)
}

type TagRepository interface {
	GetTagsByArticle(article *model.Article) ([]model.Tag, error)
	Delete(tid uint) error
	Create(tag string) (*model.Tag, error)
	List() ([]model.Tag, error)
	Update(tag *model.Tag) (*model.Tag, error)
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
