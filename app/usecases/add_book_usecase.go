package usecases

import (
	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/repositories"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/google/uuid"
)

type AddBookUsecase interface {
	Add(req request.BookRequest, LoggedUserId uuid.UUID) (*response.BookResponse, *errs.AppError)
}

type DefaultAddBookUsecase struct {
	repo        repositories.BookRepository
	userUsecase UserUsecase
}

func NewAddBookUsecase(repo repositories.BookRepository, userUsecase UserUsecase) AddBookUsecase {
	return DefaultAddBookUsecase{repo, userUsecase}
}

func (u DefaultAddBookUsecase) Add(req request.BookRequest, loggedUserId uuid.UUID) (*response.BookResponse, *errs.AppError) {
	user, err := u.userUsecase.FindUserByID(loggedUserId)
	if err != nil {
		return nil, err
	}
	book := req.ToBook(user.ID)

	createdBook, err := u.repo.CreateBook(book)
	if err != nil {
		return nil, err
	}

	response := response.BookResponse{}.ToResponse(*createdBook, *user)

	return response, nil
}
