package usecases

import (
	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/repositories"
	"github.com/Jonss/book-lending/infra/errs"
	"github.com/google/uuid"
)

type UserUsecase interface {
	Create(req request.UserRequest) (*response.UserResponse, *errs.AppError)
	FindUserByID(externalId uuid.UUID) (*response.UserResponse, *errs.AppError)
}

type DefaultUserUsecase struct {
	repo repositories.UserRepository
}

func NewCreateUserUseCase(repo repositories.UserRepository) UserUsecase {
	return DefaultUserUsecase{repo}
}

func (u DefaultUserUsecase) Create(req request.UserRequest) (*response.UserResponse, *errs.AppError) {
	user := req.ToUser()
	userCreated, err := u.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	response := response.UserResponse{}.FromUser(*userCreated)
	return &response, nil
}

func (u DefaultUserUsecase) FindUserByID(externalId uuid.UUID) (*response.UserResponse, *errs.AppError) {
	user, err := u.repo.FindUserByExternalId(externalId)
	if err != nil {
		return nil, err
	}
	response := response.UserResponse{}.FromUser(*user)
	return &response, nil
}
