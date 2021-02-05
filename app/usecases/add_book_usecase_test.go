package usecases

import (
	"net/http"
	"testing"

	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type BookRepositoryMock struct {
	mock.Mock
}

type FindUserUsecaseMock struct {
	mock.Mock
}

func (m *BookRepositoryMock) CreateBook(book models.Book) (*models.Book, *errs.AppError) {
	args := m.Called(book)
	result := args.Get(0).(*models.Book)
	return result, args.Get(1).(*errs.AppError)
}

func (m *BookRepositoryMock) BookExists(book models.Book) bool {
	args := m.Called(book)
	return args.Bool(0)
}

func (m *BookRepositoryMock) FindBookByTitleAndOwnerId(title string, userID int64) (*models.Book, *errs.AppError) {
	args := m.Called(title, userID)
	result := args.Get(0).(*models.Book)
	return result, args.Get(1).(*errs.AppError)
}

func (m *FindUserUsecaseMock) FindUserByID(externalId uuid.UUID) (*response.UserResponse, *errs.AppError) {
	args := m.Called(externalId)
	result := args.Get(0).(*response.UserResponse)
	return result, args.Get(1).(*errs.AppError)
}

func (m *FindUserUsecaseMock) FindUserByEmail(email string) (*response.UserResponse, *errs.AppError) {
	args := m.Called(email)
	result := args.Get(0).(*response.UserResponse)
	return result, args.Get(1).(*errs.AppError)
}

func TestAddBookSuccess(t *testing.T) {
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)
	bookStatusRepoMock := new(BookStatusRepositoryMock)

	externalId := uuid.New()
	expectedUser := userResponseStub(1)
	expectedBook := bookModelStub()

	bookRepoMock.On("BookExists", expectedBook).Return(false)
	findUserUsecaseMock.On("FindUserByID", externalId).Return(&expectedUser, (*errs.AppError)(nil))
	bookRepoMock.On("CreateBook", expectedBook).Return(&expectedBook, (*errs.AppError)(nil))
	bookStatusRepoMock.On("AddStatus", expectedBook, expectedUser.ID, "IDLE").Return(&expectedBook, (*errs.AppError)(nil))

	usecase := NewAddBookUsecase(bookRepoMock, findUserUsecaseMock, bookStatusRepoMock)

	result, err := usecase.Add(bookRequestStub(), externalId)

	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)
	bookStatusRepoMock.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, "Arthur C. Clarke", result.Author)
	assert.Equal(t, "O fim da infância", result.Title)
	assert.Equal(t, "Jupiter Stein", result.Owner.FullName)
	assert.Equal(t, "jupiter.stein@gmail.com", result.Owner.Email)
	assert.Equal(t, int64(1), result.Owner.ID)
	assert.NotNil(t, result.CreatedAt)
}

func TestAddBookErrorBookAlreadyInUserCollection(t *testing.T) {
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)
	bookStatusRepoMock := new(BookStatusRepositoryMock)

	externalId := uuid.New()
	expectedUser := userResponseStub(1)
	expectedBook := bookModelStub()

	bookRepoMock.On("BookExists", expectedBook).Return(true)
	findUserUsecaseMock.On("FindUserByID", externalId).Return(&expectedUser, (*errs.AppError)(nil))

	usecase := NewAddBookUsecase(bookRepoMock, findUserUsecaseMock, bookStatusRepoMock)

	result, err := usecase.Add(bookRequestStub(), externalId)

	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)
	bookStatusRepoMock.AssertExpectations(t)

	assert.Equal(t, 422, err.Code)
	assert.Equal(t, "Book [O fim da infância] by [Arthur C. Clarke] already in user 1 collection", err.Message)
	assert.Nil(t, result)
}

func TestAddBookErrorOnPersistence(t *testing.T) {
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)
	bookStatusRepoMock := new(BookStatusRepositoryMock)

	externalId := uuid.New()
	expectedUser := userResponseStub(1)
	expectedBook := bookModelStub()

	bookRepoMock.On("BookExists", expectedBook).Return(false)
	findUserUsecaseMock.On("FindUserByID", externalId).Return(&expectedUser, (*errs.AppError)(nil))
	bookRepoMock.On("CreateBook", expectedBook).Return((*models.Book)(nil), errs.NewError("Error on persistence", http.StatusInternalServerError))

	usecase := NewAddBookUsecase(bookRepoMock, findUserUsecaseMock, bookStatusRepoMock)

	result, err := usecase.Add(bookRequestStub(), externalId)

	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)
	bookStatusRepoMock.AssertExpectations(t)

	assert.Equal(t, 500, err.Code)
	assert.Equal(t, "Error on persistence", err.Message)
	assert.Nil(t, result)
}

func bookRequestStub() request.BookRequest {
	return request.BookRequest{
		Title:  "O fim da infância",
		Author: "Arthur C. Clarke",
	}
}

func bookModelStub() models.Book {
	return models.Book{
		Title:   "O fim da infância",
		Author:  "Arthur C. Clarke",
		OwnerID: 1,
		Slug:    "o-fim-da-infancia-1",
	}
}
