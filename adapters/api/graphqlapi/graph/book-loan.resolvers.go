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

	response, err := r.createUserUsecase.Create(request)
	if err != nil {
		return nil, errors.New(err.Message)
	}

	user := &model.User{
		Name:      response.FullName,
		Email:     response.Email,
		ID:        response.LoggedUserId.String(),
		CreatedAt: response.CreatedAt,
	}

	return user, nil
}

func (r *mutationResolver) AddBookToMyCollection(ctx context.Context, loggedUserID string, input model.AddBookInput) (*model.Book, error) {
	userID := uuid.MustParse(loggedUserID)

	request := request.BookRequest{
		Title: input.Title,
		Pages: input.Pages,
	}

	response, err := r.addBookUsecase.Add(request, userID)
	if err != nil {
		return nil, errors.New(err.Message)
	}

	book := &model.Book{
		ID:        response.ExternalID,
		Title:     response.Title,
		Pages:     response.Pages,
		CreatedAt: response.CreatedAt.String(),
	}
	return book, nil
}

func (r *mutationResolver) LendBook(ctx context.Context, loggedUserID string, input model.LendBookInput) (*model.BookLoan, error) {
	userID := uuid.MustParse(loggedUserID)

	request := request.LendBookRequest{
		BookID:       input.BookID,
		UserToLendID: input.ToUserID,
	}

	response, err := r.lendBookUseCase.Lend(request, userID)
	if err != nil {
		return nil, errors.New(err.Message)
	}

	bookLoan := &model.BookLoan{
		Book: &model.Book{
			ID:        response.BookResponse.ExternalID,
			Title:     response.BookResponse.Title,
			CreatedAt: response.BookResponse.CreatedAt.String(),
			Pages:     response.BookResponse.Pages,
		},
		FromUser: response.ToUserID,
		ToUser:   response.ToUserID,
		LentAt:   response.LentAt,
	}

	return bookLoan, nil
}

func (r *mutationResolver) ReturnBook(ctx context.Context, loggedUserID string, bookID string) (*model.BookLoan, error) {
	userID := uuid.MustParse(loggedUserID)

	response, err := r.returnBookUseCase.Return(bookID, userID)
	if err != nil {
		return nil, errors.New(err.Message)
	}

	bookLoan := &model.BookLoan{
		Book: &model.Book{
			ID:        response.BookResponse.ExternalID,
			Title:     response.BookResponse.Title,
			Pages:     response.BookResponse.Pages,
			CreatedAt: response.BookResponse.CreatedAt.String(),
		},
		FromUser:   response.FromUserID,
		ToUser:     response.ToUserID,
		ReturnedAt: response.BookResponse.CreatedAt.String(),
	}
	return bookLoan, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	logger.Info(fmt.Sprintf("UUID %s", id))
	loggedUserID := uuid.MustParse(id)

	response, err := r.Resolver.findUserUsecase.FindUserByID(loggedUserID)
	if err != nil {
		return nil, errors.New(err.Message)
	}

	user := &model.User{
		Name:      response.FullName,
		Email:     response.Email,
		ID:        response.LoggedUserId.String(),
		CreatedAt: response.CreatedAt,
	}
	// todo incluir collection

	return user, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
