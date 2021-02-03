package app

import (
	"fmt"

	"github.com/Jonss/book-lending/adapters/storages/pg"
	"github.com/Jonss/book-lending/adapters/storages/pg/repository"
	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/app/usecases"
)

func Start() {
	dbClient := pg.GetDbClient()
	pg.MigratePSQL(dbClient)

	userRepository := repository.NewUserRepositoryDB(dbClient)
	createUserUsecase := usecases.NewCreateUserUseCase(userRepository)
	findUserUsecase := usecases.NewFindUserUseCase(userRepository)

	bookRepository := repository.NewBookRepositoryDB(dbClient)
	addBookUseCase := usecases.NewAddBookUsecase(bookRepository, findUserUsecase)

	userRequest := request.UserRequest{Email: "jupiter.stein@gmail.com", FullName: "Júpiter Stein"}

	bookRequest := request.BookRequest{Title: "Os demônios", Author: "Fiodor Dostoievski"}

	user, err := createUserUsecase.Create(userRequest)

	fmt.Println(err)
	fmt.Println(user)

	book, err := addBookUseCase.Add(bookRequest, user.LoggedUserId)

	fmt.Println(book)
	fmt.Println(err)
}
