package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Jonss/book-lending/adapters/api/graphqlapi"
	"github.com/Jonss/book-lending/adapters/api/graphqlapi/graph"
	"github.com/Jonss/book-lending/adapters/api/rest"
	"github.com/Jonss/book-lending/adapters/storages/pg"
	"github.com/Jonss/book-lending/adapters/storages/pg/repository"
	"github.com/Jonss/book-lending/app/usecases"
	"github.com/Jonss/book-lending/infra/logger"
	"github.com/gorilla/mux"
)

func Start() {
	dbClient := pg.GetDbClient()
	pg.MigratePSQL(dbClient)

	router := mux.NewRouter()

	// dependency
	userRepository := repository.NewUserRepositoryDB(dbClient)
	bookRepository := repository.NewBookRepositoryDB(dbClient)
	bookStatusRepository := repository.NewBookStatusRepositoryDb(dbClient)

	findBookUsecase := usecases.NewFindBooksUsecase(bookRepository)
	createUserUsecase := usecases.NewCreateUserUseCase(userRepository)
	findUserUsecase := usecases.NewFindUserUseCase(userRepository, findBookUsecase)
	bookStatusUsecase := usecases.NewAddBookUsecase(bookRepository, findUserUsecase, bookStatusRepository)
	lendBookUseCase := usecases.NewLendBookUsecase(bookStatusRepository, bookRepository, findUserUsecase)
	returnBookUsecase := usecases.NewReturnBookUsecase(bookStatusRepository, findUserUsecase)

	// rest
	uh := rest.NewUserRestHandler(findUserUsecase)
	router.HandleFunc("/users/{logged_user_id}", uh.GetUser).Methods(http.MethodGet)

	// graphql
	resolver := graph.NewGraphqlResolver(
		findUserUsecase,
		createUserUsecase,
		bookStatusUsecase,
		lendBookUseCase,
		returnBookUsecase,
	)

	router.Handle("/", graphqlapi.Playground())
	srv := graphqlapi.GraphlSrv(resolver)
	router.Handle("/query", srv)

	port := os.Getenv("APP_PORT")
	logger.Info(fmt.Sprintf("Book lending is running on port %s", port))
	http.ListenAndServe(port, router)
}
