package repository

import (
	"database/sql"

	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/Jonss/book-lending/infra/logger"
)

type BookRepositoryDb struct {
	client *sql.DB
}

func NewBookRepositoryDB(client *sql.DB) BookRepositoryDb {
	return BookRepositoryDb{client}
}

func (r BookRepositoryDb) CreateBook(book models.Book) (*models.Book, *errs.AppError) {
	sql := `INSERT INTO books (title, author, owner_id, created_at) VALUES ($1, $2, $3, $4)`

	_, err := r.client.Exec(sql, book.Title, book.Author, book.OwnerID, book.CreatedAt)
	if err != nil {
		logger.Info("Error while creating new book: " + err.Error())
		return nil, errs.NewError("Unexpected error from database", 500)
	}
	return &book, nil
}
