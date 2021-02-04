package repositories

import (
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
)

type BookStatusRepository interface {
	AddStatus(models.Book, int64, string) (*models.Book, *errs.AppError)
	VerifyStatus(models.Book) *errs.AppError
}