package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"

	null "github.com/volatiletech/null/v8"

	"go-template/daos"
	"go-template/gqlmodels"
	"go-template/models"
	"go-template/pkg/utl/cnvrttogql"
	"go-template/pkg/utl/convert"
	"go-template/pkg/utl/resultwrapper"
)

// CreateAuthor is the resolver for the createAuthor field.
func (r *mutationResolver) CreateAuthor(ctx context.Context, input gqlmodels.AuthorCreateInput) (*gqlmodels.Author, error) {
	hashedPassword := r.Sec.Hash(input.Password)
	newAuthor, err := daos.CreateAuthor(ctx, models.Author{
		Email:     input.Email,
		Password:  hashedPassword,
		FirstName: input.FirstName,
		LastName:  null.StringFrom(input.LastName),
	})

	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "author insertion")
	}
	return cnvrttogql.AuthorToGraphQlAuthor(newAuthor), nil
}

// UpdateAuthor is the resolver for the updateAuthor field.
func (r *mutationResolver) UpdateAuthor(ctx context.Context, input gqlmodels.AuthorUpdateInput) (*gqlmodels.Author, error) {
	// Fetch the author
	author, err := daos.FindAuthorByID(ctx, convert.StringToInt(input.ID))
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "find author by id")
	}

	// Update Author with with with given input, if they are non null
	if input.FirstName != nil {
		author.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		author.LastName = null.StringFrom(*input.LastName)
	}

	// Update Author in the DB
	updatedAuthor, err := daos.UpdateAuthor(ctx, *author)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "update author")
	}
	// Return the updated Author
	return cnvrttogql.AuthorToGraphQlAuthor(updatedAuthor), nil
}

// DeleteAuthor is the resolver for the deleteAuthor field.
func (r *mutationResolver) DeleteAuthor(ctx context.Context, input gqlmodels.AuthorDeleteInput) (*gqlmodels.Author, error) {
	// Fetch the author
	author, err := daos.FindAuthorByID(ctx, convert.StringToInt(input.ID))
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "find author by id")
	}

	// Delete the author
	deletedCount, err := daos.DeleteAuthor(ctx, *author)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "delete author")
	}
	if deletedCount == 0 {
		return nil, resultwrapper.ResolverWrapperFromMessage(500, "no author deleted")
	}

	return cnvrttogql.AuthorToGraphQlAuthor(*author), nil
}
