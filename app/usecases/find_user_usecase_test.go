package usecases

import (
	"net/http"
	"testing"

	"github.com/Jonss/book-lending/adapters/util"
	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

type FindBooksUsecaseMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) FindUserByExternalId(externalId uuid.UUID) (*models.User, *errs.AppError) {
	args := m.Called(externalId)
	result := args.Get(0).(*models.User)
	return result, args.Get(1).(*errs.AppError)
}

func (m *FindBooksUsecaseMock) FindBooksByUserID(id int64) ([]response.BookResponse, *errs.AppError) {
	args := m.Called(id)
	result := args.Get(0).([]response.BookResponse)
	return result, args.Get(1).(*errs.AppError)
}

func TestFindUserWithSuccess(t *testing.T) {
	userRepositoryMock := new(UserRepositoryMock)
	findBookUseCaseMock := new(FindBooksUsecaseMock)

	expected := userModelStub()
	externalId := uuid.New()

	userRepositoryMock.On("FindUserByExternalId", externalId).Return(&expected, (*errs.AppError)(nil))
	findBookUseCaseMock.On("FindBooksByUserID", mock.Anything).Return(booksResponse(expected.ID), (*errs.AppError)(nil))

	usecase := NewFindUserUseCase(userRepositoryMock, findBookUseCaseMock)

	result, err := usecase.FindUserByID(externalId)

	userRepositoryMock.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, "jupiter.stein@gmail.com", result.Email)
	assert.Equal(t, "Jupiter Stein", result.FullName)
	assert.NotNil(t, result.LoggedUserId)
	assert.NotNil(t, result.Books)
	assert.NotNil(t, 3, len(result.Books))
}

func TestFindUserWithError(t *testing.T) {
	userRepositoryMock := new(UserRepositoryMock)
	findBookUseCaseMock := new(FindBooksUsecaseMock)

	externalId := uuid.New()

	userRepositoryMock.On("FindUserByExternalId", externalId).Return((*models.User)(nil), errs.NewError("User not found", http.StatusNotFound))

	usecase := NewFindUserUseCase(userRepositoryMock, findBookUseCaseMock)

	result, err := usecase.FindUserByID(externalId)

	userRepositoryMock.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, "User not found", err.Message)
	assert.Equal(t, 404, err.Code)
}

func booksResponse(userID int64) []response.BookResponse {
	return []response.BookResponse{
		{
			Title:      "Os demônios",
			Pages:      704,
			ExternalID: util.Slug("Os demônios", userID),
		},
		{
			Title:      "Memórias do subsolo",
			Pages:      152,
			ExternalID: util.Slug("Memórias do subsolo", userID),
		},
		{
			Title:      "Pais e filhos",
			Pages:      344,
			ExternalID: util.Slug("Pais e filhos", userID),
		},
	}
}
