package usecases

import (
	"fmt"
	"net/http"

	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/repositories"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/Jonss/book-lending/infra/logger"
	"github.com/google/uuid"
)

type LendBookUsecase interface {
	Lend(request.LendBookRequest, uuid.UUID) (*response.BookLoanResponse, *errs.AppError)
}

type DefaultLendBookUsecase struct {
	bookStatusRepo  repositories.BookStatusRepository
	bookRepo        repositories.BookRepository
	findUserUsecase FindUserUsecase
}

func NewLendBookUsecase(bookStatusRepo repositories.BookStatusRepository, bookRepo repositories.BookRepository, findUserUsecase FindUserUsecase) LendBookUsecase {
	return DefaultLendBookUsecase{bookStatusRepo, bookRepo, findUserUsecase}
}

func (u DefaultLendBookUsecase) Lend(req request.LendBookRequest, ownerID uuid.UUID) (*response.BookLoanResponse, *errs.AppError) {
	userID := uuid.MustParse(req.UserToLendID)
	uToLent, err := u.findUserUsecase.FindUserByID(userID)
	if err != nil {
		logger.Error(err.Message)
		message := fmt.Sprintf("User with loggedUserId [%s] not found. Cant' lend book [%s]", req.UserToLendID, req.BookID)
		logger.Error(message)
		err.Message = message
		return nil, err
	}

	uOwner, err := u.findUserUsecase.FindUserByID(ownerID)
	if err != nil {
		return nil, err
	}

	book, err := u.bookRepo.FindBookBySlug(req.BookID)
	if err != nil {
		logger.Error(err.Message)
		return nil, err
	}

	if uOwner.ID == uToLent.ID {
		logger.Error(fmt.Sprintf("Can't lend book to owner. [Book: %s. Owner: %s]", book.Title, uOwner.FullName))
		return nil, errs.NewError("Can't lend your own book", 422)
	}

	status, err := u.bookStatusRepo.VerifyStatus(*book)
	if err != nil {
		logger.Warn(fmt.Sprintf("Can't lend book to owner. [Book: %s. Owner: %s]", book.Title, uOwner.FullName))
		return nil, err
	}

	toStatus := "LENT"
	if *status == toStatus {
		return nil, errs.NewError(fmt.Sprintf("Book %s is not IDLE to be lent. Current status is %s", book.Title, *status), http.StatusUnprocessableEntity)
	}

	bookStatus, err := u.bookStatusRepo.AddStatus(*book, uToLent.ID, toStatus)
	if err != nil {
		return nil, err
	}

	fromUser := uOwner.LoggedUserId.String()
	toUser := uToLent.LoggedUserId.String()
	bookLoanResponse := response.ToBookLoanResponse(*bookStatus, fromUser, toUser, toStatus, "", bookStatus.CreatedAt.String())

	return &bookLoanResponse, nil
}
