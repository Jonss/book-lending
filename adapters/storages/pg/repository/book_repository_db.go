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
	response, err := r.client.Exec(sql, book.Title, book.Author, book.OwnerID, book.CreatedAt)
	if err != nil {
		logger.Info("Error while creating new book: " + err.Error())
		return nil, errs.NewError("Unexpected error from database", 500)
	}

	id, err := response.LastInsertId()
	if err != nil {
		logger.Info("Error while getting book id: " + err.Error())
		return nil, errs.NewError("Unexpected error from database", 500)
	}

	book.ID = id

	return &book, nil
}

func (r BookRepositoryDb) BookExists(book models.Book) bool {
	sql := `SELECT COUNT(*) FROM books WHERE owner_id = $1 AND title = $2 AND author = $3`

	var counter int
	err := r.client.QueryRow(sql, book.OwnerID, book.Title, book.Author).Scan(&counter)
	if err != nil {
		return true
	}
	return counter > 0
}

func (r BookRepositoryDb) FindBookByTitleAndOwnerId(title string, ownerID int64) (*models.Book, *errs.AppError) {
	logger.Info(fmt.Sprintf("Search book [%s]", title))

	sql := `SELECT (id, title, author, owner_id, created_at) from books where owner_id = $1 AND title = $2`

	rows, err := r.client.Query(sql, ownerID, title)
	if err != nil {
		logger.Error("Error fetching book: " + err.Error())
		return nil, errs.NewError("book not found", 404)
	}

	for rows.Next() {
		var b models.Book
		rows.Scan(&b.ID, &b.Title, &b.Author, &b.OwnerID, &b.CreatedAt)
		logger.Info(fmt.Sprintf("Book found: [Title: %s - OwnerID: %d]", b.Title, b.OwnerID))
		return &b, nil
	}

	return nil, errs.NewError("book not found", 404)
}
