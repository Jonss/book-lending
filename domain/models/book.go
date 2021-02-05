package models

import "time"

type Book struct {
	ID        int64
	Title     string
	Author    string
	OwnerID   int64
	CreatedAt time.Time
	Slug      string
}
