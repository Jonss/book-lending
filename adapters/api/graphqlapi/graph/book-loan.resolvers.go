package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jonss/book-lending/adapters/api/graphqlapi/graph/generated"
	"github.com/Jonss/book-lending/adapters/api/graphqlapi/graph/model"
	"github.com/Jonss/book-lending/app/dto/request"
	"github.com/Jonss/book-lending/infra/logger"
	"github.com/google/uuid"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	request := request.UserRequest{
		Email:    input.Email,
		FullName: input.Name,
	}

	user, err := r.createUserUsecase.Create(request)
	if err != nil {
		return nil, errors.New(err.Message)
	}

	gphUser := &model.User{
		Name:      user.FullName,
		Email:     user.Email,
		ID:        user.LoggedUserId.String(),
		CreatedAt: user.CreatedAt,
	}

	return gphUser, nil
}

func (r *mutationResolver) AddBookToMyCollection(ctx context.Context, loggedUserID string, input model.AddBookInput) (*model.Book, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) LendBook(ctx context.Context, loggedUserID string, input model.LendBookInput) (*model.BookLoan, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ReturnBook(ctx context.Context, loggedUserID string, bookID string) (*model.BookLoan, error) {

	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	logger.Info(fmt.Sprintf("UUID %s", id))
	loggedUserID := uuid.MustParse(id)

	user, err := r.Resolver.findUserUsecase.FindUserByID(loggedUserID)
	if err != nil {
		return nil, errors.New(err.Message)
	}

	gphUser := &model.User{
		Name:      user.FullName,
		Email:     user.Email,
		ID:        user.LoggedUserId.String(),
		CreatedAt: user.CreatedAt,
	}

	return gphUser, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
