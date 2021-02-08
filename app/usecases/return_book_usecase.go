package usecases

import (
	"fmt"
	"net/http"

	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/repositories"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/Jonss/book-lending/infra/logger"
	"github.com/google/uuid"
)

type ReturnBookUsecase interface {
	Return(string, uuid.UUID) (*response.BookLoanResponse, *errs.AppError)
}

type DefaultReturnBookUsecase struct {
	bookStatusRepo  repositories.BookStatusRepository
	findUserUsecase FindUserUsecase
}

func NewReturnBookUsecase(bookStatusRepo repositories.BookStatusRepository, FindUserUsecase FindUserUsecase) ReturnBookUsecase {
	return DefaultReturnBookUsecase{bookStatusRepo, FindUserUsecase}
}

func (u DefaultReturnBookUsecase) Return(slug string, userUUID uuid.UUID) (*response.BookLoanResponse, *errs.AppError) {
	user, err := u.findUserUsecase.FindUserByID(userUUID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error searching user %s to return book [%s]. ", userUUID, slug) + err.Message)
		return nil, err
	}

	bookStatus, err := u.bookStatusRepo.FindStatusBySlug(slug)
	if err != nil {
		logger.Error(fmt.Sprintf("Error searching book %s to return book. ", slug) + err.Message)
		return nil, err
	}

	toStatus := "IDLE"
	if bookStatus.Status == toStatus {
		return nil, errs.NewError(fmt.Sprintf("Book %s is not LENT to be returned. Current status is %s", bookStatus.Book.Title, bookStatus.Status), http.StatusUnprocessableEntity)
	}

	if bookStatus.BearerUserID != user.ID || bookStatus.Status != "LENT" {
		logger.Warn(fmt.Sprintf("User returning book is not the same with book. Returning userID: %d, User with book: %d", user.ID, bookStatus.BearerUserID))
		return nil, errs.NewError("can't return book", http.StatusUnprocessableEntity)
	}

	status := "IDLE"
	u.bookStatusRepo.AddStatus(*bookStatus.Book, bookStatus.Book.OwnerID, status)
	logger.Info(fmt.Sprintf("Book returned. %s", slug))

	fromUser := bookStatus.Book.Owner.LoggedUserId.String()
	toUser := user.LoggedUserId.String()
	bookLoanResponse := response.ToBookLoanResponse(*bookStatus, fromUser, toUser, status, bookStatus.CreatedAt.String(), "")

	return &bookLoanResponse, nil
}
