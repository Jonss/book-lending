package response

import (
	"time"

	"github.com/Jonss/book-lending/adapters/util"
	"github.com/Jonss/book-lending/domain/models"
)

type BookResponse struct {
	Title     string       `json:"title"`
	Author    string       `json:"author"`
	Owner     UserResponse `json:"owner"`
	CreatedAt time.Time    `json:"created_at"`
	Status    string       `json:"status"`
	Slug      string       `json:"title_slug"`
}

func (r BookResponse) ToResponse(book models.Book, userResponse UserResponse, status string) *BookResponse {
	return &BookResponse{
		Title:     book.Title,
		Author:    book.Author,
		Owner:     userResponse,
		CreatedAt: book.CreatedAt,
		Status:    status,
		Slug:      util.Slug(book.Title, book.ID),
	}
}
