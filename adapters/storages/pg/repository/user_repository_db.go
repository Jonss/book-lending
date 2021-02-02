package repository

import (
	"database/sql"
	"fmt"

	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/Jonss/book-lending/infra/logger"
	"github.com/google/uuid"
)

const duplicatedUserMessage string = `pq: duplicate key value violates unique constraint "users_email_key"`

type UserRepositoryDb struct {
	client *sql.DB
}

func NewUserRepositoryDB(client *sql.DB) UserRepositoryDb {
	return UserRepositoryDb{client}
}

func (r UserRepositoryDb) CreateUser(user models.User) (*models.User, *errs.AppError) {
	sql := `INSERT INTO users (full_name, external_id, email) VALUES ($1, $2, $3)`

	user.LoggedUserId = uuid.New()
	_, err := r.client.Exec(sql, user.FullName, user.LoggedUserId, user.Email)
	if err != nil {
		if duplicatedUserMessage == err.Error() {
			logger.Warn("User already exists: " + err.Error())
			return nil, errs.NewError(fmt.Sprintf("user with email %s already exists", user.Email), 422)
		}
		logger.Info("Error while creating new account: " + err.Error())
		return nil, errs.NewError("Unexpected error from database", 500)
	}
	return &user, nil
}
