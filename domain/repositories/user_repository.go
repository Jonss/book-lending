package repositories

import (
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(models.User) (*models.User, *errs.AppError)
	FindUserByExternalId(uuid.UUID) (*models.User, *errs.AppError)
}
