package usecases

import (
	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/app/dto/response"
	"github.com/Jonss/book-lending/domain/repositories"
	"github.com/Jonss/book-lending/infra/errs"
)

type CreateUserUsecase interface {
	Create(req request.UserRequest) (*response.UserResponse, *errs.AppError)
}

type DefaultCreateUserUsecase struct {
	repo repositories.UserRepository
}

func NewCreateUserUseCase(repo repositories.UserRepository) CreateUserUsecase {
	return DefaultCreateUserUsecase{repo}
}

func (u DefaultCreateUserUsecase) Create(req request.UserRequest) (*response.UserResponse, *errs.AppError) {
	user := req.ToUser()
	userCreated, err := u.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	response := response.FromUser(*userCreated)
	return &response, nil
}
