package graph

import "github.com/Jonss/book-lending/app/usecases"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	findUserUsecase   usecases.FindUserUsecase
	createUserUsecase usecases.CreateUserUsecase
	addBookUsecase    usecases.AddBookUsecase
}

func NewGraphqlResolver(findUserUsecase usecases.FindUserUsecase, createUserUsecase usecases.CreateUserUsecase, addBookUsecase usecases.AddBookUsecase) Resolver {
	return Resolver{findUserUsecase, createUserUsecase, addBookUsecase}
}
