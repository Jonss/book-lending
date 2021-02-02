package usecases

import (
	"testing"

	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) CreateUser(user models.User) (*models.User, *errs.AppError) {
	args := m.Called(user)
	result := args.Get(0).(*models.User)
	return result, args.Get(1).(*errs.AppError)
}

func (m *UserRepositoryMock) FindUserByExternalId(externalId uuid.UUID) (*models.User, *errs.AppError) {
	args := m.Called(externalId)
	result := args.Get(0).(*models.User)
	return result, args.Get(1).(*errs.AppError)
}

func TestCreateUserWithSuccess(t *testing.T) {
	repo := new(UserRepositoryMock)

	expected := userModel()

	repo.On("CreateUser", expected).Return(&expected, (*errs.AppError)(nil))

	usecase := NewCreateUserUseCase(repo)

	result, err := usecase.Create(userRequest())

	repo.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, "jupiter.stein@gmail.com", result.Email)
	assert.Equal(t, "Jupiter Stein", result.FullName)
	assert.NotNil(t, result.LoggedUserId)
}

func TestCreateUserWithError(t *testing.T) {
	repo := new(UserRepositoryMock)

	expected := userModel()

	repo.On("CreateUser", expected).Return((*models.User)(nil), errs.NewError("Error", 500))

	usecase := NewCreateUserUseCase(repo)

	result, err := usecase.Create(userRequest())

	repo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, "Error", err.Message)
	assert.Equal(t, 500, err.Code)
}

func userRequest() request.UserRequest {
	return request.UserRequest{
		Email:    "jupiter.stein@gmail.com",
		FullName: "Jupiter Stein",
	}
}

func userModel() models.User {
	return models.User{
		Email:        "jupiter.stein@gmail.com",
		FullName:     "Jupiter Stein",
		LoggedUserId: uuid.UUID{},
	}
}
