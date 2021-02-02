package repositories

import (
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
)

type UserRepository interface {
	CreateUser(models.User) (*models.User, *errs.AppError)
}
