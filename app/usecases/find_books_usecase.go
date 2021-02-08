package usecases

import (
	"fmt"

	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/repositories"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/Jonss/book-lending/infra/logger"
)

type FindBooksUsecase interface {
	FindBooksByUserID(int64) ([]response.BookResponse, *errs.AppError)
}

type DefaultFindBooksUsecase struct {
	repo repositories.BookRepository
}

func NewFindBooksUsecase(repo repositories.BookRepository) FindBooksUsecase {
	return DefaultFindBooksUsecase{repo}
}

func (u DefaultFindBooksUsecase) FindBooksByUserID(userID int64) ([]response.BookResponse, *errs.AppError) {

	books, err := u.repo.FindBooksByOwner(userID)
	if err != nil {
		logger.Error(fmt.Sprintf("Fetching books from user %d", userID))
		return []response.BookResponse{}, err
	}

	logger.Info(fmt.Sprintf("user %d has %d books", userID, len(books)))

	var booksResponse []response.BookResponse
	for _, v := range books {
		br := response.ToBookResponse(v)
		booksResponse = append(booksResponse, br)
	}

	return booksResponse, nil
}
