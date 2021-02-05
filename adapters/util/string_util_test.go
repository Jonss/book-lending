package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldTransformStringToSlug(t *testing.T) {
	title := "Deset pražských dnů, 17.-27. listopad 1989: Dokumentace (Czech Edition)"
	id := int64(1)

	slug := Slug(title, id)
	assert.Equal(t, "deset-prazskych-dnu-17-27-listopad-1989-dokumentace-czech-edition-1", slug)
}
