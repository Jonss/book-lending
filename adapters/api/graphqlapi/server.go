package graphqlapi

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Jonss/book-lending/adapters/api/graphqlapi/graph"
	"github.com/Jonss/book-lending/adapters/api/graphqlapi/graph/generated"
)

func Playground() http.HandlerFunc {
	return playground.Handler("GraphQL playground", "/query")
}

func GraphlSrv(r graph.Resolver) *handler.Server {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &r}))
}
