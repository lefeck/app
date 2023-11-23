package service

import (
	"app/common/request"
	"app/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type RBACService interface {
	List() ([]model.Role, error)
	ListResources() ([]model.Resource, error)
	Create(role *model.Role) (*model.Role, error)
	CreateResource(resource *model.Resource) (*model.Resource, error)
	CreateResources(resources []model.Resource, conds ...clause.Expression) error
	GetRoleByID(id string) (*model.Role, error)
	GetResource(id string) (*model.Resource, error)
	GetRoleByName(name string) (*model.Role, error)
	Update(id string, role *model.Role) (*model.Role, error)
	Delete(id string) error
	DeleteResource(id string) error
	ListOperations() ([]model.Operation, error)
}

type UserService interface {
	List(pageSize int, pageNum int) (int, []interface{})
	Create(user *model.User) (*model.User, error)
	Get(string) (*model.User, error)
	CreateOAuthUser(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(string) error
	FindAll(userlist []model.User) ([]model.User, error)
	Validate(*model.User) error
	//Auth(*AuthUser) (*model.User, error)
	Login(param request.Login) (*model.User, error)
	Export(data *[]model.User, headerName []string, filename string, c *gin.Context) error
}

type ArticleService interface {
	List() ([]model.Article, error)
	Create(*model.User, *model.Article) (*model.Article, error)
	Get(user *model.User, id string) (*model.Article, error)
	Update(id string, post *model.Article) (*model.Article, error)
	Delete(id string) error
}

type LikeService interface {
}

type TagService interface {
	GetTagsByArticle(id string) ([]model.Tag, error)
	Create(tag string) (*model.Tag, error)
	Delete(tid uint)
	Update(tag *model.Tag, id string) (*model.Tag, error)
}

type CategoryService interface {
}
