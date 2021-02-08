package usecases

import (
	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/domain/repositories"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/google/uuid"
)

type FindUserUsecase interface {
	FindUserByID(externalId uuid.UUID) (*response.UserResponse, *errs.AppError)
}

type DefaultFindUserUsecase struct {
	repo            repositories.UserRepository
	findBookUsecase FindBooksUsecase
}

func NewFindUserUseCase(repo repositories.UserRepository, findBookUsecase FindBooksUsecase) FindUserUsecase {
	return DefaultFindUserUsecase{repo, findBookUsecase}
}

func (u DefaultFindUserUsecase) FindUserByID(externalId uuid.UUID) (*response.UserResponse, *errs.AppError) {
	user, err := u.repo.FindUserByExternalId(externalId)
	if err != nil {
		return nil, err
	}

	books, err := u.findBookUsecase.FindBooksByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	return handleResponse(user, err, books)
}

func handleResponse(user *models.User, err *errs.AppError, books []response.BookResponse) (*response.UserResponse, *errs.AppError) {
	if err != nil {
		return nil, err
	}
	response := response.FromUser(*user)
	response.Books = books
	return &response, nil
}
