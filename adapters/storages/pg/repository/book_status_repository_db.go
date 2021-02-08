package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/Jonss/book-lending/infra/logger"
	"github.com/google/uuid"
)

type BookStatusRepositoryDb struct {
	client *sql.DB
}

func NewBookStatusRepositoryDb(client *sql.DB) BookStatusRepositoryDb {
	return BookStatusRepositoryDb{client}
}

func (r BookStatusRepositoryDb) AddStatus(book models.Book, bearerID int64, status string) (*models.BookStatus, *errs.AppError) {
	logger.Info(fmt.Sprintf("Add book status. Book: %s - %d. ToUserID: %d Status: [%s]", book.Title, book.ID, bearerID, status))

	createdAt := time.Now()

	if status == "IDLE" {
		updatePreviousStatus(r.client, book.ID)
	}

	sql := "INSERT INTO books_status(book_id, bearer_user_id, status, created_at) VALUES ($1, $2, $3, $4)"

	_, err := r.client.Exec(sql, book.ID, bearerID, status, createdAt)
	if err != nil {
		logger.Info("Error while creating books_status account: " + err.Error())
		return nil, errs.NewError("Unexpected error from database", http.StatusInternalServerError)
	}

	bookStatus := &models.BookStatus{
		Book:         &book,
		Status:       status,
		BearerUserID: bearerID,
		CreatedAt:    createdAt,
	}

	return bookStatus, nil
}

func (r BookStatusRepositoryDb) VerifyStatus(book models.Book) (*string, *errs.AppError) {
	logger.Info(fmt.Sprintf("Verifying book status. Book: %s. Owner: %d", book.Slug, book.OwnerID))

	sql := `SELECT status FROM books_status where book_id = $1 ORDER BY id DESC LIMIT 1`
	row := r.client.QueryRow(sql, book.ID)

	var status string
	errScan := row.Scan(&status)
	if errScan != nil {
		logger.Error(fmt.Sprintf("Error getting book status: %s", errScan.Error()))
		return nil, errs.NewError("Error getting book status", 500)
	}

	return &status, nil
}

func (r BookStatusRepositoryDb) FindStatusBySlug(slug string) (*models.BookStatus, *errs.AppError) {
	logger.Info(fmt.Sprintf("Search book slug [%s]", slug))

	sql := `
	SELECT b.id,
		b.title,
		b.owner_id,
		b.created_at,
		bs.id,
		bs.status,
		bs.bearer_user_id,
		u.external_id
		from books_status bs
		inner join books b
		on b.id = bs.book_id
		INNER JOIN users u
		ON b.id = u.id
		where b.slug = $1
		order by bs.id desc limit 1;`

	row := r.client.QueryRow(sql, slug)

	var bookID, ownerID, bookStatusID, bearerUserID int64
	var title, status, userExternalID string
	var createdAt time.Time
	errScan := row.Scan(&bookID, &title, &ownerID, &createdAt, &bookStatusID, &status, &bearerUserID, &userExternalID)

	if errScan != nil {
		logger.Error("Error getting book status " + errScan.Error())
		return nil, errs.NewError("Book status not found", http.StatusNotFound)
	}

	bookStatus := models.BookStatus{
		Status:       status,
		BearerUserID: bearerUserID,
		Book: &models.Book{
			ID:        bookID,
			Title:     title,
			OwnerID:   ownerID,
			CreatedAt: createdAt,
			Owner: models.User{
				ID:           ownerID,
				LoggedUserId: uuid.MustParse(userExternalID),
			},
		},
	}

	return &bookStatus, nil
}

func updatePreviousStatus(client *sql.DB, bookID int64) {
	sql := "select id from books_status where book_id = $1 order by id desc limit 1"

	row := client.QueryRow(sql, bookID)

	var bookStatusID int
	err := row.Scan(&bookStatusID)
	if err != nil {
		logger.Error("error searching book_status")
	}

	sql = "update books_status set returned_at = $1 where id = $2"

	_, err = client.Exec(sql, time.Now(), bookStatusID)
	if err != nil {
		logger.Error(fmt.Sprintf("error updating latest book_status with id %d, %s", bookStatusID, err.Error()))
	}
}
