package request

import (
	"github.com/Jonss/book-lending/adapters/util"
	"github.com/Jonss/book-lending/domain/models"
)

type BookRequest struct {
	Title string `json:"title"`
	Pages int    `author:"pages"`
}

func (r BookRequest) ToBook(userID int64) models.Book {
	return models.Book{
		Title:   r.Title,
		OwnerID: userID,
		Pages:   r.Pages,
		Slug:    util.Slug(r.Title, userID),
	}
}
