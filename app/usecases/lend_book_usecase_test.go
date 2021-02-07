package usecases

import (
	"fmt"
	"testing"

	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type BookStatusRepositoryMock struct {
	mock.Mock
}

func (m *BookStatusRepositoryMock) AddStatus(book models.Book, userLenderID int64, status string) (*models.Book, *errs.AppError) {
	args := m.Called(book, userLenderID, status)
	result := args.Get(0).(*models.Book)
	return result, args.Get(1).(*errs.AppError)
}

func (m *BookStatusRepositoryMock) VerifyStatus(book models.Book) *errs.AppError {
	args := m.Called(book)
	return args.Get(0).(*errs.AppError)
}

func (m *BookStatusRepositoryMock) FindStatusBySlug(slug string) (*models.BookStatus, *errs.AppError) {
	args := m.Called(slug)
	result := args.Get(0).(*models.BookStatus)
	return result, args.Get(1).(*errs.AppError)
}

func TestLendBookToUserWithSuccess(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	ownerUUID := uuid.New()
	lenderUUID := uuid.New()
	expectedUser := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerUUID)
	expectedBook := bookModelStub()
	expectedOwner := userResponseToBuildStub(2, "machado@gmail.com", "Machado", lenderUUID)
	lendBookRequest := lendBookRequestStub()

	findUserUsecaseMock.On("FindUserByEmail", "machado@gmail.com").Return(&expectedUser, (*errs.AppError)(nil))
	findUserUsecaseMock.On("FindUserByID", ownerUUID).Return(&expectedOwner, (*errs.AppError)(nil))
	bookRepoMock.On("FindBookByTitleAndOwnerId", lendBookRequest.Title, expectedOwner.ID).Return(&expectedBook, (*errs.AppError)(nil))
	bookStatusRepoMock.On("VerifyStatus", expectedBook).Return((*errs.AppError)(nil))
	bookStatusRepoMock.On("AddStatus", expectedBook, expectedUser.ID, "LENT").Return(&expectedBook, (*errs.AppError)(nil))

	usecase := NewDefaultLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

	result, err := usecase.Lend(lendBookRequest, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, "O fim da infância", result.Title)
	assert.Equal(t, 299, result.Pages)
	assert.Equal(t, int64(1), result.OwnerID)
	assert.NotNil(t, result.CreatedAt)
}

func TestLendBookToUserWithErrorWhenLenderDoesNotExists(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	ownerId := uuid.New()
	expectedUser := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerId)
	lendBookRequest := lendBookRequestStub()

	findUserUsecaseMock.On("FindUserByEmail", "machado@gmail.com").Return((*response.UserResponse)(nil), errs.NewError("user not found", 404))

	usecase := NewDefaultLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

	result, err := usecase.Lend(lendBookRequest, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "User with email [machado@gmail.com] not found. Cant' lend book [O fim da infância]", err.Message)
}

func TestLendBookToUserWithErrorOwnerIsNotFound(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	ownerUUID := uuid.New()
	expectedUser := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerUUID)
	lendBookRequest := lendBookRequestStub()

	findUserUsecaseMock.On("FindUserByEmail", "machado@gmail.com").Return(&expectedUser, (*errs.AppError)(nil))
	findUserUsecaseMock.On("FindUserByID", ownerUUID).Return((*response.UserResponse)(nil), errs.NewError("user not found", 404))

	usecase := NewDefaultLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

	result, err := usecase.Lend(lendBookRequest, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "user not found", err.Message)
}

func TestLendBookToUserWithErrorWhenBookDoesNotExists(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	ownerUUID := uuid.New()
	lenderUUID := uuid.New()
	expectedUser := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerUUID)
	expectedOwner := userResponseToBuildStub(2, "machado@gmail.com", "Machado", lenderUUID)
	lendBookRequest := lendBookRequestStub()

	findUserUsecaseMock.On("FindUserByEmail", "machado@gmail.com").Return(&expectedUser, (*errs.AppError)(nil))
	findUserUsecaseMock.On("FindUserByID", ownerUUID).Return(&expectedOwner, (*errs.AppError)(nil))
	bookRepoMock.On("FindBookByTitleAndOwnerId", lendBookRequest.Title, expectedOwner.ID).Return((*models.Book)(nil), errs.NewError("book not found", 404))

	usecase := NewDefaultLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

	result, err := usecase.Lend(lendBookRequest, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "book not found", err.Message)
}

func TestLendBookToUserWithErrorWhenBookStatusIsNotIdle(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	ownerUUID := uuid.New()
	lenderUUID := uuid.New()
	expectedUser := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerUUID)
	expectedBook := bookModelStub()
	expectedOwner := userResponseToBuildStub(2, "machado@gmail.com", "Machado", lenderUUID)
	lendBookRequest := lendBookRequestStub()

	findUserUsecaseMock.On("FindUserByEmail", "machado@gmail.com").Return(&expectedUser, (*errs.AppError)(nil))
	findUserUsecaseMock.On("FindUserByID", ownerUUID).Return(&expectedOwner, (*errs.AppError)(nil))
	bookRepoMock.On("FindBookByTitleAndOwnerId", lendBookRequest.Title, expectedOwner.ID).Return(&expectedBook, (*errs.AppError)(nil))
	bookStatusRepoMock.On("VerifyStatus", expectedBook).Return(errs.NewError(fmt.Sprintf("Book is not IDLE. Current status is LENT"), 422))

	usecase := NewDefaultLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

	result, err := usecase.Lend(lendBookRequest, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, 422, err.Code)
	assert.Equal(t, "Book is not IDLE. Current status is LENT", err.Message)
}

func lendBookRequestStub() request.LendBookRequest {
	return request.LendBookRequest{
		Title:           "O fim da infância",
		UserToLendEmail: "machado@gmail.com",
	}
}
