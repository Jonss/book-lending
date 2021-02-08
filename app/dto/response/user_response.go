package response

import (
	"github.com/Jonss/book-lending/domain/models"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID           int64          `json:",omitempty"`
	LoggedUserId uuid.UUID      `json:"logged_user_id"`
	FullName     string         `json:"full_name"`
	Email        string         `json:"email"`
	CreatedAt    string         `json:"created_at"`
	Books        []BookResponse `json:"books"`
}

func FromUser(user models.User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		LoggedUserId: user.LoggedUserId,
		FullName:     user.FullName,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt.Local().String(),
	}
}
