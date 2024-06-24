package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"
	"fmt"

	null "github.com/volatiletech/null/v8"

	"go-template/daos"
	"go-template/gqlmodels"
	"go-template/internal/middleware/auth"
	"go-template/pkg/utl/convert"
	"go-template/pkg/utl/resultwrapper"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*gqlmodels.LoginResponse, error) {
	u, err := daos.FindUserByUserName(username, ctx)
	if err != nil {
		return nil, err
	}
	if !u.Password.Valid || (!r.Sec.HashMatchesPassword(u.Password.String, password)) {
		return nil, fmt.Errorf("username or password does not exist ")
	}

	if !u.Active.Valid || (!u.Active.Bool) {
		return nil, resultwrapper.ErrUnauthorized
	}

	token, err := r.JWTService.GenerateToken(u)
	if err != nil {
		return nil, resultwrapper.ErrUnauthorized
	}

	refreshToken := r.Sec.Token(token)
	u.Token = null.StringFrom(refreshToken)
	_, err = daos.UpdateUser(*u, ctx)
	if err != nil {
		return nil, err
	}

	return &gqlmodels.LoginResponse{Token: token, RefreshToken: refreshToken}, nil
}

// AuthorLogin is the resolver for the authorLogin field.
func (r *mutationResolver) AuthorLogin(ctx context.Context, email string, password string) (*gqlmodels.LoginResponse, error) {
	author, err := daos.FindAuthorByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("username or password does not exist haha")
	}

	if !r.Sec.HashMatchesPassword(author.Password, password) {
		return nil, fmt.Errorf("username or password does not exist ")
	}

	token, err := r.JWTService.GenerateTokenForAuthor(author)
	if err != nil {
		return nil, resultwrapper.ErrUnauthorized
	}

	refreshToken := r.Sec.Token(token)

	return &gqlmodels.LoginResponse{Token: token, RefreshToken: refreshToken}, nil
}

// ChangePassword is the resolver for the changePassword field.
func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword string, newPassword string) (*gqlmodels.ChangePasswordResponse, error) {
	userID := auth.UserIDFromContext(ctx)
	u, err := daos.FindUserByID(userID, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}

	if !r.Sec.HashMatchesPassword(convert.NullDotStringToString(u.Password), oldPassword) {
		return nil, fmt.Errorf("incorrect old password")
	}

	if !r.Sec.Password(newPassword,
		convert.NullDotStringToString(u.FirstName),
		convert.NullDotStringToString(u.LastName),
		convert.NullDotStringToString(u.Username),
		convert.NullDotStringToString(u.Email)) {
		return nil, fmt.Errorf("inr.Secure password")
	}

	u.Password = null.StringFrom(r.Sec.Hash(newPassword))
	_, err = daos.UpdateUser(*u, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "new information")
	}
	return &gqlmodels.ChangePasswordResponse{Ok: true}, err
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (*gqlmodels.RefreshTokenResponse, error) {
	user, err := daos.FindUserByToken(token, ctx)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "token")
	}
	resp, err := r.JWTService.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	return &gqlmodels.RefreshTokenResponse{Token: resp}, nil
}

// Mutation returns gqlmodels.MutationResolver implementation.
func (r *Resolver) Mutation() gqlmodels.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
