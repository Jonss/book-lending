package graph

import "github.com/Jonss/book-lending/app/usecases"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	findUserUsecase   usecases.FindUserUsecase
	createUserUsecase usecases.CreateUserUsecase
	addBookUsecase    usecases.AddBookUsecase
	lendBookUseCase   usecases.LendBookUsecase
	returnBookUseCase usecases.ReturnBookUsecase
}

func NewGraphqlResolver(findUserUsecase usecases.FindUserUsecase,
	createUserUsecase usecases.CreateUserUsecase,
	addBookUsecase usecases.AddBookUsecase,
	lendBookUseCase usecases.LendBookUsecase,
	returnBookUseCase usecases.ReturnBookUsecase,
) Resolver {
	return Resolver{findUserUsecase, createUserUsecase, addBookUsecase, lendBookUseCase, returnBookUseCase}
}
