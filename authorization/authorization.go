package authorization

import (
	"app/repository"
)

var store repository.Repository

func InitAuthorization(repository repository.Repository) error {
	store = repository
	return nil
}
