package request

import (
	"github.com/Jonss/book-lending/domain/models"
	"github.com/google/uuid"
)

type UserRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

func (r UserRequest) ToUser() models.User {
	return models.User{
		LoggedUserId: uuid.New(),
		Email:        r.Email,
		FullName:     r.FullName,
	}
}
