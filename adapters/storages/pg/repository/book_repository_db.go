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

func (r BookRepositoryDb) FindBookBySlug(slug string) (*models.Book, *errs.AppError) {
	logger.Info(fmt.Sprintf("Search book [%s]", slug))

	sql := `SELECT id, title, owner_id, created_at, pages, slug from books where slug = $1`

	rows, err := r.client.Query(sql, slug)
	if err != nil {
		logger.Error("Error fetching book: " + err.Error())
		return nil, errs.NewError("book not found", http.StatusNotFound)
	}

	for rows.Next() {
		var b models.Book
		rows.Scan(&b.ID, &b.Title, &b.OwnerID, &b.CreatedAt, &b.Pages, &b.Slug)
		logger.Info(fmt.Sprintf("Book found: [Title: %s - OwnerID: %d]", b.Title, b.OwnerID))
		return &b, nil
	}

	logger.Warn(fmt.Sprintf("Book not found: [Slug: %s]", slug))
	return nil, errs.NewError("book not found", http.StatusNotFound)
}

func (r BookRepositoryDb) FindBooksByOwner(userID int64) ([]models.Book, *errs.AppError) {
	sql := `SELECT id, title, owner_id, created_at, pages, slug from books where owner_id = $1`

	rows, err := r.client.Query(sql, userID)
	if err != nil {
		logger.Error("Error fetching book: " + err.Error())
		return nil, errs.NewError("book not found", http.StatusNotFound)
	}

	var books []models.Book
	for rows.Next() {
		var b models.Book
		rows.Scan(&b.ID, &b.Title, &b.OwnerID, &b.CreatedAt, &b.Pages, &b.Slug)
		logger.Info(fmt.Sprintf("Book found: [Title: %s - OwnerID: %d]", b.Title, b.OwnerID))
		books = append(books, b)
	}
	return books, nil
}
