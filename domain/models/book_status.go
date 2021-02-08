package models

import "time"

type BookStatus struct {
	Status       string
	BearerUserID int64
	Book         *Book
	CreatedAt    time.Time
}
