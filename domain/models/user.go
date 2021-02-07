package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           int64
	LoggedUserId uuid.UUID `db:"external_id"`
	FullName     string
	Email        string
	CreatedAt    time.Time
}
