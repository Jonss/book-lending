package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookRequestToBook(t *testing.T) {
	bookRequest := BookRequest{
		Title:  "A morte de Ivan Ilitch",
		Author: "Lev Tolstoi",
	}

	book := bookRequest.ToBook(2)

	assert.Equal(t, bookRequest.Author, book.Author)
	assert.Equal(t, bookRequest.Title, book.Title)
	assert.Equal(t, int64(2), book.OwnerID)
	assert.Equal(t, "a-morte-de-ivan-ilitch-2", book.Slug)
}
