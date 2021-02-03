package request

import (
	"github.com/Jonss/book-lending/domain/models"
)

type BookRequest struct {
	Title  string `json:"title"`
	Author string `author:"author"`
}

func (r BookRequest) ToBook(userId int64) models.Book {
	return models.Book{
		Author:  r.Author,
		Title:   r.Title,
		OwnerID: userId,
	}
}
