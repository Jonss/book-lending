package models

import "github.com/google/uuid"

type User struct {
	ID           int64
	LoggedUserId uuid.UUID `db:"external_id"`
	FullName     string
	Email        string
}
