package response

import (
	"time"

	"github.com/Jonss/book-lending/domain/models"
)

type BookResponse struct {
	Title     string       `json:"title"`
	Author    string       `json:"author"`
	Owner     UserResponse `json:"owner"`
	CreatedAt time.Time    `json:"created_at"`
}

func (r BookResponse) ToResponse(book models.Book, userResponse UserResponse) *BookResponse {
	return &BookResponse{
		Title:     book.Title,
		Author:    book.Author,
		Owner:     userResponse,
		CreatedAt: book.CreatedAt,
	}
}
