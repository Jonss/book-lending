package models

import "time"

type Book struct {
	ID        int64
	Title     string
	OwnerID   int64
	CreatedAt time.Time
	Pages     int
	Slug      string
}
