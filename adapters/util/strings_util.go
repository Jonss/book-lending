package util

import (
	"fmt"
	"strconv"
	"strings"

	slugify "github.com/metal3d/go-slugify"
)

func Slug(title string, id int64) string {
	strID := strconv.FormatInt(id, 10)
	return strings.ToLower(fmt.Sprintf("%s-%s", slugify.Marshal(title), strID))
}
