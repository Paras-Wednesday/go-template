package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"
	"fmt"
	"go-template/gqlmodels"
)

// PostByID is the resolver for the postByID field.
func (r *queryResolver) PostByID(ctx context.Context, id string) (*gqlmodels.Post, error) {
	panic(fmt.Errorf("not implemented: PostByID - postByID"))
}

// AllPostByAuthor is the resolver for the allPostByAuthor field.
func (r *queryResolver) AllPostByAuthor(ctx context.Context, authorID string) (*gqlmodels.PostsPayload, error) {
	panic(fmt.Errorf("not implemented: AllPostByAuthor - allPostByAuthor"))
}
