package models

type BookStatus struct {
	Status       string
	BearerUserID int64
	Book         Book
}
