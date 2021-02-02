package request

import (
	"github.com/Jonss/book-lending/domain/models"
)

type UserRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

func (r UserRequest) ToUser() models.User {
	return models.User{
		Email:    r.Email,
		FullName: r.FullName,
	}
}
