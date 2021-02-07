package response

import (
	"time"

	"github.com/Jonss/book-lending/domain/models"
)

type BookResponse struct {
	Title     string       `json:"title"`
	Owner     UserResponse `json:"owner"`
	CreatedAt time.Time    `json:"created_at"`
	Status    string       `json:"status"`
	Pages     int          `json:"pages"`
}

func (r BookResponse) ToResponse(book models.Book, userResponse UserResponse, status string) *BookResponse {
	return &BookResponse{
		Title:     book.Title,
		Owner:     userResponse,
		CreatedAt: book.CreatedAt,
		Status:    status,
		Pages:     book.Pages,
	}
}
