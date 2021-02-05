package usecases

import (
	"fmt"
	"net/http"

	"github.com/Jonss/book-lending/domain/repositories"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/Jonss/book-lending/infra/logger"
	"github.com/google/uuid"
)

type ReturnBookUsecase interface {
	Return(string, uuid.UUID) *errs.AppError
}

type DefaultReturnBookUsecase struct {
	bookStatusRepo  repositories.BookStatusRepository
	findUserUsecase FindUserUsecase
}

func NewDefaultReturnBookUsecase(bookStatusRepo repositories.BookStatusRepository, FindUserUsecase FindUserUsecase) ReturnBookUsecase {
	return DefaultReturnBookUsecase{bookStatusRepo, FindUserUsecase}
}

func (u DefaultReturnBookUsecase) Return(slug string, userUUID uuid.UUID) *errs.AppError {
	user, err := u.findUserUsecase.FindUserByID(userUUID)
	if err != nil {
		logger.Error(fmt.Sprintf("Error searching user %s to return book [%s]. ", userUUID, slug) + err.Message)
		return err
	}

	book, err := u.bookStatusRepo.FindStatusBySlug(slug)
	if err != nil {
		logger.Error(fmt.Sprintf("Error searching book %s to return book. ", slug) + err.Message)
		return err
	}

	if book.BearerUserID != user.ID || book.Status != "LENT" {
		logger.Warn(fmt.Sprintf("User returning book is not the same with book. Returning userID: %d, User with book: %d", user.ID, book.BearerUserID))
		return errs.NewError("can't return book", http.StatusUnprocessableEntity)
	}

	u.bookStatusRepo.AddStatus(book.Book, book.Book.OwnerID, "IDLE")
	logger.Info(fmt.Sprintf("Book returned. %s", slug))
	return nil
}
