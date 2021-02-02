package usecases

import (
	"testing"

	"github.com/Jonss/book-lending/domain/models"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) FindUserByExternalId(externalId uuid.UUID) (*models.User, *errs.AppError) {
	args := m.Called(externalId)
	result := args.Get(0).(*models.User)
	return result, args.Get(1).(*errs.AppError)
}

func TestFindUserWithSuccess(t *testing.T) {
	repo := new(UserRepositoryMock)

	expected := userModel()

	externalId := uuid.New()

	repo.On("FindUserByExternalId", externalId).Return(&expected, (*errs.AppError)(nil))

	usecase := NewFindUserUseCase(repo)

	result, err := usecase.FindUserByID(externalId)

	repo.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, "jupiter.stein@gmail.com", result.Email)
	assert.Equal(t, "Jupiter Stein", result.FullName)
	assert.NotNil(t, result.LoggedUserId)
}

func TestFindUserWithError(t *testing.T) {
	repo := new(UserRepositoryMock)

	externalId := uuid.New()

	repo.On("FindUserByExternalId", externalId).Return((*models.User)(nil), errs.NewError("User not found", 404))

	usecase := NewFindUserUseCase(repo)

	result, err := usecase.FindUserByID(externalId)

	repo.AssertExpectations(t)

	assert.Nil(t, result)
	assert.Equal(t, "User not found", err.Message)
	assert.Equal(t, 404, err.Code)
}
