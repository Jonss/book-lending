package request

import (
	"github.com/Jonss/book-lending/adapters/util"
	"github.com/Jonss/book-lending/domain/models"
)

type BookRequest struct {
	Title  string `json:"title"`
	Author string `author:"author"`
}

func (r BookRequest) ToBook(userID int64) models.Book {
	return models.Book{
		Author:  r.Author,
		Title:   r.Title,
		OwnerID: userID,
		Slug:    util.Slug(r.Title, userID),
	}
}
