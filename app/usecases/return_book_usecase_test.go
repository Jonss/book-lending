package usecases

import (
	"net/http"
	"testing"

	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReturnBookSuccess(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	expectedUser := userResponseStub(1)
	expectedBook := bookModelStub()
	expectedBookStatus := bookStatusStub(expectedBook, 1, "LENT")

	findUserUsecaseMock.On("FindUserByID", mock.Anything).Return(&expectedUser, (*errs.AppError)(nil))
	bookStatusRepoMock.On("FindStatusBySlug", expectedBook.Slug).Return(&expectedBookStatus, (*errs.AppError)(nil))
	bookStatusRepoMock.On("AddStatus", expectedBook, expectedUser.ID, "IDLE").Return(&expectedBook, (*errs.AppError)(nil))

	usecase := NewDefaultReturnBookUsecase(bookStatusRepoMock, findUserUsecaseMock)

	err := usecase.Return(expectedBook.Slug, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.Nil(t, err)
}

func TestReturnBookErrorWhenUserNotFound(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	expectedUser := userResponseStub(1)
	expectedBook := bookModelStub()

	findUserUsecaseMock.On("FindUserByID", mock.Anything).Return((*response.UserResponse)(nil), errs.NewError("user not found", http.StatusNotFound))

	usecase := NewDefaultReturnBookUsecase(bookStatusRepoMock, findUserUsecaseMock)

	err := usecase.Return(expectedBook.Slug, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.Equal(t, "user not found", err.Message)
	assert.Equal(t, 404, err.Code)
}

func TestReturnBookErrorWhenBearerIDIsDifferentOfUserID(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	expectedUser := userResponseStub(1)
	expectedBook := bookModelStub()
	expectedBookStatus := bookStatusStub(expectedBook, 1, "LENT")

	findUserUsecaseMock.On("FindUserByID", mock.Anything).Return(&expectedUser, (*errs.AppError)(nil))
	bookStatusRepoMock.On("FindStatusBySlug", expectedBook.Slug).Return(&expectedBookStatus, (*errs.AppError)(nil))
	bookStatusRepoMock.On("AddStatus", expectedBook, expectedUser.ID, "IDLE").Return(&expectedBook, (*errs.AppError)(nil))

	usecase := NewDefaultReturnBookUsecase(bookStatusRepoMock, findUserUsecaseMock)

	err := usecase.Return(expectedBook.Slug, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.Nil(t, err)
}

func TestReturnBookWithErrorWhenBookIsNotLent(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	expectedUser := userResponseStub(1)
	expectedBook := bookModelStub()
	expectedBookStatus := bookStatusStub(expectedBook, 1, "IDLE")

	findUserUsecaseMock.On("FindUserByID", mock.Anything).Return(&expectedUser, (*errs.AppError)(nil))
	bookStatusRepoMock.On("FindStatusBySlug", expectedBook.Slug).Return(&expectedBookStatus, (*errs.AppError)(nil))

	usecase := NewDefaultReturnBookUsecase(bookStatusRepoMock, findUserUsecaseMock)

	err := usecase.Return(expectedBook.Slug, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.Equal(t, "can't return book", err.Message)
	assert.Equal(t, 422, err.Code)
}

func TestReturnBookWithErrorWhenBookBearerIDIsDifferentOfUserID(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	expectedUser := userResponseStub(1)
	expectedBook := bookModelStub()
	expectedBookStatus := bookStatusStub(expectedBook, 2, "LENT")

	findUserUsecaseMock.On("FindUserByID", mock.Anything).Return(&expectedUser, (*errs.AppError)(nil))
	bookStatusRepoMock.On("FindStatusBySlug", expectedBook.Slug).Return(&expectedBookStatus, (*errs.AppError)(nil))

	usecase := NewDefaultReturnBookUsecase(bookStatusRepoMock, findUserUsecaseMock)

	err := usecase.Return(expectedBook.Slug, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.Equal(t, "can't return book", err.Message)
	assert.Equal(t, 422, err.Code)
}

func bookStatusStub(book models.Book, id int64, status string) models.BookStatus {
	return models.BookStatus{
		Book:         book,
		BearerUserID: id,
		Status:       status,
	}
}
