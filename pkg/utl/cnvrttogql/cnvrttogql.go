package cnvrttogql

import (
	"context"
	"strconv"

	"github.com/volatiletech/sqlboiler/v4/boil"

	graphql "go-template/gqlmodels"
	"go-template/internal/constants"
	"go-template/models"
	"go-template/pkg/utl/convert"
)

// UsersToGraphQlUsers converts array of type models.User into array of pointer type graphql.User
func UsersToGraphQlUsers(u models.UserSlice, count int) []*graphql.User {
	var r []*graphql.User
	for _, e := range u {
		r = append(r, UserToGraphQlUser(e, count))
	}
	return r
}

// UserToGraphQlUser converts type models.User into pointer type graphql.User
func UserToGraphQlUser(u *models.User, count int) *graphql.User {
	count++
	if u == nil {
		return nil
	}
	var role *models.Role
	if count <= constants.MaxDepth {
		u.L.LoadRole(context.Background(), boil.GetContextDB(), true, u, nil) //nolint:errcheck
		if u.R != nil {
			role = u.R.Role
		}
	}

	return &graphql.User{
		ID:        strconv.Itoa(u.ID),
		FirstName: convert.NullDotStringToPointerString(u.FirstName),
		LastName:  convert.NullDotStringToPointerString(u.LastName),
		Username:  convert.NullDotStringToPointerString(u.Username),
		Email:     convert.NullDotStringToPointerString(u.Email),
		Mobile:    convert.NullDotStringToPointerString(u.Mobile),
		Address:   convert.NullDotStringToPointerString(u.Address),
		Active:    convert.NullDotBoolToPointerBool(u.Active),
		Role:      RoleToGraphqlRole(role, count),
	}
}

func RoleToGraphqlRole(r *models.Role, count int) *graphql.Role {
	count++
	if r == nil {
		return nil
	}
	var users models.UserSlice
	if count <= constants.MaxDepth {
		r.L.LoadUsers(context.Background(), boil.GetContextDB(), true, r, nil) //nolint:errcheck
		if r.R != nil {
			users = r.R.Users
		}
	}

	return &graphql.Role{
		ID:          strconv.Itoa(r.ID),
		AccessLevel: r.AccessLevel,
		Name:        r.Name,
		UpdatedAt:   convert.NullDotTimeToPointerInt(r.UpdatedAt),
		CreatedAt:   convert.NullDotTimeToPointerInt(r.CreatedAt),
		Users:       UsersToGraphQlUsers(users, count),
	}
}

func AuthorToGraphQlAuthor(a models.Author) *graphql.Author {
	return &graphql.Author{
		ID:        strconv.Itoa(a.ID),
		FirstName: a.FirstName,
		LastName:  a.LastName.String,
		Email:     a.Email,
		CreatedAt: convert.NullDotTimeToPointerInt(a.CreatedAt),
		UpdatedAt: convert.NullDotTimeToPointerInt(a.UpdatedAt),
	}
}

func AuthorsToGraphQlAuthorsPayload(authors models.AuthorSlice, total int64) *graphql.AuthorsPayload {
	result := graphql.AuthorsPayload{
		Authors: make([]*graphql.Author, 0, len(authors)),
		Total:   int(total),
	}

	for i := range authors {
		result.Authors = append(result.Authors, AuthorToGraphQlAuthor(*authors[i]))
	}
	return &result
}

func PostToGraphqlPost(post *models.Post) *graphql.Post {
	if post == nil {
		return nil
	}

	return &graphql.Post{
		ID:        strconv.Itoa(post.ID),
		AuthorID:  strconv.Itoa(post.AuthorID),
		Content:   post.Content,
		CreatedAt: convert.NullDotTimeToPointerInt(post.CreatedAt),
		UpdatedAt: convert.NullDotTimeToPointerInt(post.UpdatedAt),
	}
}

func PostsToGraphqlPostsPayload(postSlice models.PostSlice, count int64) *graphql.PostsPayload {
	result := graphql.PostsPayload{
		Posts: make([]*graphql.Post, 0, len(postSlice)),
		Total: int(count),
	}
	for i := range postSlice {
		result.Posts = append(result.Posts, PostToGraphqlPost(postSlice[i]))
	}
	return &result
}
