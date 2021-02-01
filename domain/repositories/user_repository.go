package repositories

import (
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra"
)

type UserRepository interface {
	CreateUser(models.User) (*models.User, *infra.AppError)
}
