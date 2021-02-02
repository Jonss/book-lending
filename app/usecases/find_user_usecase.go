package usecases

import (
	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/repositories"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/google/uuid"
)

type FindUserUsecase interface {
	FindUserByID(externalId uuid.UUID) (*response.UserResponse, *errs.AppError)
}

type DefaultFindUserUsecase struct {
	repo repositories.UserRepository
}

func NewFindUserUseCase(repo repositories.UserRepository) FindUserUsecase {
	return DefaultFindUserUsecase{repo}
}

func (u DefaultFindUserUsecase) FindUserByID(externalId uuid.UUID) (*response.UserResponse, *errs.AppError) {
	user, err := u.repo.FindUserByExternalId(externalId)
	if err != nil {
		return nil, err
	}
	response := response.UserResponse{}.FromUser(*user)
	return &response, nil
}
