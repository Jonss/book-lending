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
		return nil, errs.NewError(fmt.Sprintf("Book [%s] already in user %d collection", book.Title, book.OwnerID), 422)
	}

	createdBook, err := u.repo.CreateBook(book)
	if err != nil {
		return nil, err
	}

	var status string = "IDLE"
	b, err := u.bookStatusRepo.AddStatus(*createdBook, user.ID, status)
	if err != nil {
		return nil, err
	}

	bookResponse := response.BookResponse{
		Title:      createdBook.Title,
		ExternalID: createdBook.Slug,
		Owner: response.UserResponse{
			LoggedUserId: b.Book.Owner.LoggedUserId,
			FullName:     b.Book.Owner.FullName,
			Email:        b.Book.Owner.Email,
		},
		CreatedAt: createdBook.CreatedAt,
		Status:    b.Status,
		Pages:     b.Book.Pages,
	}

	return &bookResponse, nil
}
