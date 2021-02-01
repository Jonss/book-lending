package main

import (
	"fmt"
	"os"

	"github.com/Jonss/book-lending/adapters/storages/pg"
	"github.com/Jonss/book-lending/adapters/storages/pg/repository"
	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/app/usecases"
	"github.com/Jonss/book-lending/infra/logger"
	_ "github.com/lib/pq"
)

func main() {
	logger.Info(fmt.Sprintf("Starting application in environment: [%s]", os.Getenv("ENV")))
	dbClient := pg.GetDbClient()
	pg.MigratePSQL(dbClient)

	userRepository := repository.NewUserRepositoryDB(dbClient)
	userUseCase := usecases.NewCreateUserUseCase(userRepository)

	request := request.UserRequest{Email: "jupiter.stein@gmail.com", FullName: "Júpiter Stein"}
	fmt.Println(request)
	user, err := userUseCase.Execute(request)

	fmt.Println(err)
	fmt.Println(user)

	fmt.Println("Ahoy World")
}
