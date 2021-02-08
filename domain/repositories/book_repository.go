package repositories

import (
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
)

type BookRepository interface {
	CreateBook(models.Book) (*models.Book, *errs.AppError)
	BookExists(models.Book) bool
	FindBookBySlug(string) (*models.Book, *errs.AppError)
	FindBooksByOwner(int64) ([]models.Book, *errs.AppError)
}
