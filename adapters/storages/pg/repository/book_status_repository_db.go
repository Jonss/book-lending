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

type BookStatusRepositoryDb struct {
	client *sql.DB
}

func NewBookStatusRepositoryDb(client *sql.DB) BookStatusRepositoryDb {
	return BookStatusRepositoryDb{client}
}

func (r BookStatusRepositoryDb) AddStatus(book models.Book, userLenderID int64, status string) (*models.Book, *errs.AppError) {
	logger.Info(fmt.Sprintf("Add book status. Book: %s - %d. Owner: %d, Lender: %d. Status: [%s]", book.Title, book.ID, book.OwnerID, userLenderID, status))

	sql := "INSERT INTO books_status(book_id, bearer_user_id, status) VALUES ($1, $2, $3)"

	_, err := r.client.Exec(sql, book.ID, userLenderID, status)
	if err != nil {
		logger.Info("Error while creating book_status account: " + err.Error())
		return nil, errs.NewError("Unexpected error from database", http.StatusInternalServerError)
	}
	return &book, nil
}

func (r BookStatusRepositoryDb) VerifyStatus(book models.Book) *errs.AppError {
	logger.Info(fmt.Sprintf("Verifying book status. Book: %s. Owner: %d", book.Title, book.OwnerID))

	sql := `SELECT status FROM book_status where id = $1 ORDER BY id DESC LIMIT 1`
	row := r.client.QueryRow(sql, book.ID)

	var status string
	errScan := row.Scan(&status)
	if errScan != nil {
		logger.Error("Error getting book status" + errScan.Error())
		return errs.NewError("Error getting book status", 500)
	}

	if status != "IDLE" {
		return errs.NewError(fmt.Sprintf("Book is not IDLE. Current status is %s", status), http.StatusUnprocessableEntity)
	}

	return nil
}

func (r BookStatusRepositoryDb) FindStatusBySlug(slug string) (*models.BookStatus, *errs.AppError) {
	logger.Info(fmt.Sprintf("Search book slug [%s]", slug))

	sql := `
		SELECT b.id, b.title, b.owner_id, b.created_at,
		bs.id, bs.status, bs.bearer_user_id
		FROM books b
		INNER JOIN books_status bs
		ON b.id = bs.book_id
		WHERE b.slug = $1`

	row := r.client.QueryRow(sql, slug)

	var bookID, ownerID, bookStatusID, bearerUserID int64
	var title, status string
	var createdAt time.Time
	errScan := row.Scan(&bookID, &title, &ownerID, &createdAt, &bookStatusID, &status, &bearerUserID)

	if errScan != nil {
		logger.Error("Error getting book status" + errScan.Error())
		return nil, errs.NewError("Book status not found", http.StatusNotFound)
	}

	bookStatus := models.BookStatus{
		Status:       status,
		BearerUserID: bearerUserID,
		Book: models.Book{
			ID:      bookID,
			Title:   title,
			OwnerID: ownerID,
		},
	}

	return &bookStatus, nil
}
