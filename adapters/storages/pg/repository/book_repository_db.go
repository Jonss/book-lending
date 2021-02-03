package repository

import (
	"database/sql"
	"fmt"
	"time"

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

	book.CreatedAt = time.Now()
	_, err := r.client.Exec(sql, book.Title, book.Author, book.OwnerID, book.CreatedAt)
	if err != nil {
		logger.Info("Error while creating new book: " + err.Error())
		return nil, errs.NewError("Unexpected error from database", 500)
	}
	return &book, nil
}

func (r BookRepositoryDb) BookExists(book models.Book) bool {
	fmt.Println(book)
	sql := `SELECT COUNT(*) FROM books WHERE owner_id = $1 AND title = $2 AND author = $3`

	var counter int
	err := r.client.QueryRow(sql, book.OwnerID, book.Title, book.Author).Scan(&counter)
	if err != nil {
		return true
	}
	return counter > 0
}
