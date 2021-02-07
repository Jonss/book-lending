package repository

import (
	"database/sql"
	"fmt"
	"net/http"
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
	sql := `INSERT INTO books (title, owner_id, created_at, pages, slug) VALUES ($1, $2, $3, $4, $5) returning id`

	var lastInsertID int64
	book.CreatedAt = time.Now()
	err := r.client.QueryRow(sql, book.Title, book.OwnerID, book.CreatedAt, book.Pages, book.Slug).Scan(&lastInsertID)
	if err != nil {
		logger.Info("Error while creating new book: " + err.Error())
		return nil, errs.NewError("Unexpected error from database", http.StatusInternalServerError)
	}

	book.ID = lastInsertID

	return &book, nil
}

func (r BookRepositoryDb) BookExists(book models.Book) bool {
	sql := `SELECT COUNT(*) FROM books WHERE owner_id = $1 AND title = $2`

	var counter int
	err := r.client.QueryRow(sql, book.OwnerID, book.Title).Scan(&counter)
	if err != nil {
		return true
	}
	return counter > 0
}

func (r BookRepositoryDb) FindBookByTitleAndOwnerId(title string, ownerID int64) (*models.Book, *errs.AppError) {
	logger.Info(fmt.Sprintf("Search book [%s]", title))

	sql := `SELECT (id, title, owner_id, created_at, pages) from books where owner_id = $1 AND title = $2`

	rows, err := r.client.Query(sql, ownerID, title)
	if err != nil {
		logger.Error("Error fetching book: " + err.Error())
		return nil, errs.NewError("book not found", http.StatusNotFound)
	}

	for rows.Next() {
		var b models.Book
		rows.Scan(&b.ID, &b.Title, &b.OwnerID, &b.CreatedAt, &b.Pages)
		logger.Info(fmt.Sprintf("Book found: [Title: %s - OwnerID: %d]", b.Title, b.OwnerID))
		return &b, nil
	}

	logger.Info(fmt.Sprintf("Book found: [Title: %s - OwnerID: %d]", title, ownerID))
	return nil, errs.NewError("book not found", http.StatusNotFound)
}
