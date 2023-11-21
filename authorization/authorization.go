package authorization

import (
	"app/model"
	"app/repository"
)

var store repository.Repository

func InitAuthorization(repository repository.Repository) error {
	store = repository
	return nil
}

func Authorize(user *model.User) {

}
