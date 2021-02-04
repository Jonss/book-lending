package usecases

import (
	"fmt"

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
	repo            repositories.BookRepository
	findUserUsecase FindUserUsecase
	bookStatusRepo  repositories.BookStatusRepository
}

func NewAddBookUsecase(repo repositories.BookRepository, usecase FindUserUsecase, bookStatusRepo repositories.BookStatusRepository) AddBookUsecase {
	return DefaultAddBookUsecase{repo, usecase, bookStatusRepo}
}

func (u DefaultAddBookUsecase) Add(req request.BookRequest, loggedUserId uuid.UUID) (*response.BookResponse, *errs.AppError) {
	user, err := u.findUserUsecase.FindUserByID(loggedUserId)
	if err != nil {
		return nil, err
	}

	book := req.ToBook(user.ID)

	exists := u.repo.BookExists(book)
	if exists {
		return nil, errs.NewError(fmt.Sprintf("Book [%s] by [%s] already in user %d collection", book.Title, book.Author, book.OwnerID), 422)
	}

	createdBook, err := u.repo.CreateBook(book)
	if err != nil {
		return nil, err
	}

	_, err = u.bookStatusRepo.AddStatus(book, user.ID, "IDLE")
	if err != nil {
		return nil, err
	}

	response := response.BookResponse{}.ToResponse(*createdBook, *user)

	return response, nil
}
