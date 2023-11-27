package service

import (
	"app/common/request"
	"app/model"
	"github.com/gin-gonic/gin"
)

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
	Create(user *model.User, pid string) error
	Delete(user *model.User, pid string) error
}

type TagService interface {
	GetTagsByArticle(id string) ([]model.Tag, error)
	Create(tag string) (*model.Tag, error)
	Delete(id string) error
	Update(tag *model.Tag, id string) (*model.Tag, error)
}

type CategoryService interface {
	GetCategories(id string) ([]model.Category, error)
	Create(category *model.Category) (*model.Category, error)
	Delete(id string) error
	Update(id string, category *model.Category) (*model.Category, error)
}

type CommentService interface {
	Add(comment *model.Comment, id string, user *model.User) (*model.Comment, error)
	Delete(id string) error
	List(aid string) ([]model.Comment, error)
}
