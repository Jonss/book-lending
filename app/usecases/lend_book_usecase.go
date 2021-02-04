package usecases

import (
	"fmt"

	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/domain/repositories"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/Jonss/book-lending/infra/logger"
	"github.com/google/uuid"
)

type LendBookUsecase interface {
	Lend(request.LendBookRequest, uuid.UUID) (*models.Book, *errs.AppError)
}

type DefaultLendBookUsecase struct {
	bookStatusRepo  repositories.BookStatusRepository
	bookRepo        repositories.BookRepository
	findUserUsecase FindUserUsecase
}

func NewDefaultLendBookUsecase(bookStatusRepo repositories.BookStatusRepository, bookRepo repositories.BookRepository, findUserUsecase FindUserUsecase) LendBookUsecase {
	return DefaultLendBookUsecase{bookStatusRepo, bookRepo, findUserUsecase}
}

func (u DefaultLendBookUsecase) Lend(req request.LendBookRequest, ownerID uuid.UUID) (*models.Book, *errs.AppError) {
	uToLend, err := u.findUserUsecase.FindUserByEmail(req.UserToLendEmail)
	if err != nil {
		logger.Error(err.Message)
		message := fmt.Sprintf("User with email [%s] not found. Cant' lend book [%s]", req.UserToLendEmail, req.Title)
		logger.Error(message)
		err.Message = message
		return nil, err
	}

	uOwner, err := u.findUserUsecase.FindUserByID(ownerID)
	if err != nil {
		return nil, err
	}

	book, err := u.bookRepo.FindBookByTitleAndOwnerId(req.Title, uOwner.ID)
	if err != nil {
		logger.Error(err.Message)
		return nil, err
	}

	if uOwner.ID == uToLend.ID {
		logger.Error(fmt.Sprintf("Can't lend book to owner. [Book: %s. Owner: %s]", book.Title, uOwner.FullName))
		return nil, errs.NewError("Can't lend your own book", 422)
	}

	err = u.bookStatusRepo.VerifyStatus(*book)
	if err != nil {
		logger.Warn(fmt.Sprintf("Can't lend book to owner. [Book: %s. Owner: %s]", book.Title, uOwner.FullName))
		return nil, err
	}

	u.bookStatusRepo.AddStatus(*book, uToLend.ID, "LENT")

	return book, nil
}
