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

func (m *BookStatusRepositoryMock) AddStatus(book models.Book, userLenderID int64, status string) (*models.BookStatus, *errs.AppError) {
	args := m.Called(book, userLenderID, status)
	result := args.Get(0).(*models.BookStatus)
	return result, args.Get(1).(*errs.AppError)
}

func (m *BookStatusRepositoryMock) VerifyStatus(book models.Book) (*string, *errs.AppError) {
	args := m.Called(book)
	result := args.Get(0).(*string)
	return result, args.Get(1).(*errs.AppError)
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

	ownerUUID := uuid.MustParse("b511ba57-85e4-499f-8c84-bce8d682d21c")
	lenderUUID := uuid.MustParse("1320480d-d88c-48bd-802c-89932970aa4b")
	expectedOwner := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerUUID)
	expectedBook := bookModelStub()
	expecteLender := userResponseToBuildStub(2, "machado@gmail.com", "Machado", lenderUUID)
	lendBookRequest := lendBookRequestStub()
	expectedBookStatus := bookStatusModelStub()

	status := "IDLE"
	findUserUsecaseMock.On("FindUserByID", lenderUUID).Return(&expecteLender, (*errs.AppError)(nil))
	findUserUsecaseMock.On("FindUserByID", ownerUUID).Return(&expectedOwner, (*errs.AppError)(nil))
	bookRepoMock.On("FindBookBySlug", lendBookRequest.BookID).Return(&expectedBook, (*errs.AppError)(nil))
	bookStatusRepoMock.On("VerifyStatus", expectedBook).Return(&status, (*errs.AppError)(nil))
	bookStatusRepoMock.On("AddStatus", expectedBook, expecteLender.ID, "LENT").Return(&expectedBookStatus, (*errs.AppError)(nil))

	usecase := NewLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

	result, err := usecase.Lend(lendBookRequest, expectedOwner.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, "O fim da infância", result.BookResponse.Title)
	assert.Equal(t, 299, result.BookResponse.Pages)
	assert.Equal(t, "1320480d-d88c-48bd-802c-89932970aa4b", result.ToUserID)
	assert.Equal(t, "b511ba57-85e4-499f-8c84-bce8d682d21c", result.FromUserID)
	assert.Equal(t, "LENT", result.BookResponse.Status)
}

func TestLendBookToUserWithErrorWhenBookIsAlreadyLent(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	ownerUUID := uuid.MustParse("b511ba57-85e4-499f-8c84-bce8d682d21c")
	lenderUUID := uuid.MustParse("1320480d-d88c-48bd-802c-89932970aa4b")
	expectedOwner := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerUUID)
	expectedBook := bookModelStub()
	expecteLender := userResponseToBuildStub(2, "machado@gmail.com", "Machado", lenderUUID)
	lendBookRequest := lendBookRequestStub()

	status := "LENT"
	findUserUsecaseMock.On("FindUserByID", lenderUUID).Return(&expecteLender, (*errs.AppError)(nil))
	findUserUsecaseMock.On("FindUserByID", ownerUUID).Return(&expectedOwner, (*errs.AppError)(nil))
	bookRepoMock.On("FindBookBySlug", lendBookRequest.BookID).Return(&expectedBook, (*errs.AppError)(nil))
	bookStatusRepoMock.On("VerifyStatus", expectedBook).Return(&status, (*errs.AppError)(nil))

	usecase := NewLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

	result, err := usecase.Lend(lendBookRequest, expectedOwner.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Book O fim da infância is not IDLE to be lent. Current status is LENT", err.Message)
	assert.Equal(t, 422, err.Code)
}

func TestLendBookToUserWithErrorWhenLenderDoesNotExists(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	ownerId := uuid.New()
	expectedUser := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerId)
	lendBookRequest := lendBookRequestStub()

	findUserUsecaseMock.On("FindUserByID", mock.Anything).Return((*response.UserResponse)(nil), errs.NewError("user not found", 404))

	usecase := NewLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

	result, err := usecase.Lend(lendBookRequest, expectedUser.LoggedUserId)

	bookStatusRepoMock.AssertExpectations(t)
	bookRepoMock.AssertExpectations(t)
	findUserUsecaseMock.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "User with loggedUserId [1320480d-d88c-48bd-802c-89932970aa4b] not found. Cant' lend book [o-fim-da-infancia-1]", err.Message)
}

func TestLendBookToUserWithErrorOwnerIsNotFound(t *testing.T) {
	bookStatusRepoMock := new(BookStatusRepositoryMock)
	bookRepoMock := new(BookRepositoryMock)
	findUserUsecaseMock := new(FindUserUsecaseMock)

	ownerUUID := uuid.MustParse("1320480d-d88c-48bd-802c-89932970aa4b")
	inexistingUUID := uuid.MustParse("ddf0e91b-54ae-451b-946c-2ed6b2f61554")
	expectedUser := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerUUID)
	lendBookRequest := lendBookRequestStub()

	findUserUsecaseMock.On("FindUserByID", expectedUser.LoggedUserId).Return(&expectedUser, (*errs.AppError)(nil))
	findUserUsecaseMock.On("FindUserByID", inexistingUUID).Return((*response.UserResponse)(nil), errs.NewError("user not found", 404))

	usecase := NewLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

	result, err := usecase.Lend(lendBookRequest, inexistingUUID)

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
	lenderUUID := uuid.MustParse("1320480d-d88c-48bd-802c-89932970aa4b")
	expectedUser := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerUUID)
	expectedOwner := userResponseToBuildStub(2, "machado@gmail.com", "Machado", lenderUUID)
	lendBookRequest := lendBookRequestStub()

	findUserUsecaseMock.On("FindUserByID", lenderUUID).Return(&expectedUser, (*errs.AppError)(nil))
	findUserUsecaseMock.On("FindUserByID", ownerUUID).Return(&expectedOwner, (*errs.AppError)(nil))
	bookRepoMock.On("FindBookBySlug", lendBookRequest.BookID).Return((*models.Book)(nil), errs.NewError("book not found", 404))

	usecase := NewLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

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
	lenderUUID := uuid.MustParse("1320480d-d88c-48bd-802c-89932970aa4b")
	expectedUser := userResponseToBuildStub(1, "jupiter.stein@gmail.com.", "Júpiter Stein", ownerUUID)
	expectedBook := bookModelStub()
	expectedOwner := userResponseToBuildStub(2, "machado@gmail.com", "Machado", lenderUUID)
	lendBookRequest := lendBookRequestStub()

	findUserUsecaseMock.On("FindUserByID", lenderUUID).Return(&expectedUser, (*errs.AppError)(nil))
	findUserUsecaseMock.On("FindUserByID", ownerUUID).Return(&expectedOwner, (*errs.AppError)(nil))
	bookRepoMock.On("FindBookBySlug", lendBookRequest.BookID).Return(&expectedBook, (*errs.AppError)(nil))
	bookStatusRepoMock.On("VerifyStatus", expectedBook).Return(errs.NewError(fmt.Sprintf("Book is not IDLE. Current status is LENT"), 422))

	usecase := NewLendBookUsecase(bookStatusRepoMock, bookRepoMock, findUserUsecaseMock)

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
		BookID:       "o-fim-da-infancia-1",
		UserToLendID: "1320480d-d88c-48bd-802c-89932970aa4b",
	}
}
