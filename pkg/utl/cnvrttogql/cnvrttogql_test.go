package cnvrttogql

import (
	graphql "go-template/gqlmodels"
	"go-template/models"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const SuccessCase = "Success"

func TestUsersToGraphQlUsers(t *testing.T) {
	type args struct {
		u models.UserSlice
	}
	tests := []struct {
		name string
		args args
		want []*graphql.User
	}{
		{
			name: SuccessCase,
			args: args{
				u: models.UserSlice{{
					ID: 1,
				}},
			},
			want: []*graphql.User{
				{
					ID: "1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UsersToGraphQlUsers(tt.args.u, 1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UsersToGraphQlUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleToGraphqlRole(t *testing.T) {
	type args struct {
		u *models.Role
	}
	tests := []struct {
		name string
		args args
		want *graphql.Role
	}{
		{
			name: SuccessCase,
			args: args{
				u: &models.Role{
					ID: 1,
				},
			},
			want: &graphql.Role{
				ID: "1",
			},
		},
		{
			name: SuccessCase,
			args: args{
				u: nil,
			},
			want: nil,
		},
	}

	db, _, err := sqlmock.New()
	if err != nil {
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
	}
	boil.SetDB(db)
	defer db.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoleToGraphqlRole(tt.args.u, 1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoleToGraphqlRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserToGraphQlUser(t *testing.T) {
	tests := []struct {
		name string
		req  *models.User
		want *graphql.User
	}{
		{
			name: SuccessCase,
			req:  nil,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserToGraphQlUser(tt.req, 0)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestAuthorToGraphqlAuthor(t *testing.T) {
	now := time.Now()
	nowMilli := int(now.UnixMilli())
	tests := []struct {
		name     string
		input    models.Author
		expected *graphql.Author
	}{
		{
			name: SuccessCase,
			input: models.Author{
				ID:        29,
				FirstName: "First",
				LastName:  null.String{},
				CreatedAt: null.TimeFrom(now),
				Email:     "some.email@domain.com",
				Password:  "This shouldn't be displayed",
			},
			expected: &graphql.Author{
				ID:        "29",
				FirstName: "First",
				LastName:  "",
				Email:     "some.email@domain.com",
				CreatedAt: &nowMilli,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, AuthorToGraphQlAuthor(test.input))
		})
	}
}

func TestAuthorsToGraphqlAuthorsPayload(t *testing.T) {
	now := time.Now()
	nowMilli := int(now.UnixMilli())
	tests := []struct {
		name       string
		inputSlice models.AuthorSlice
		inputTotal int64
		expected   *graphql.AuthorsPayload
	}{
		{
			name: SuccessCase,
			inputSlice: []*models.Author{
				{
					ID:        29,
					FirstName: "First",
					LastName:  null.String{},
					CreatedAt: null.TimeFrom(now),
					Password:  "Shouldn't be displayed at any cost",
				},
			},
			inputTotal: 2,
			expected: &graphql.AuthorsPayload{
				Authors: []*graphql.Author{
					{
						ID:        "29",
						FirstName: "First",
						LastName:  "",
						CreatedAt: &nowMilli,
					},
				},
				Total: 2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := AuthorsToGraphQlAuthorsPayload(test.inputSlice, test.inputTotal)
			assert.Equal(t, test.expected, got)
		})
	}
}

func TestPostToGraphqlPost(t *testing.T) {
	now := time.Now()
	nowMilli := int(now.UnixMilli())
	tests := []struct {
		name     string
		input    *models.Post
		expected *graphql.Post
	}{
		{
			name:     "nil post should return nil",
			input:    nil,
			expected: nil,
		},
		{
			name: SuccessCase,
			input: &models.Post{
				ID:        1,
				AuthorID:  2,
				Content:   "HI this is post",
				CreatedAt: null.TimeFrom(now),
			},
			expected: &graphql.Post{
				ID:        "1",
				AuthorID:  "2",
				Content:   "HI this is post",
				CreatedAt: &nowMilli,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(
				t, test.expected, PostToGraphqlPost(test.input))
		})
	}
}

func TestPostsToGraphqlPostsPayload(t *testing.T) {
	tests := []struct {
		name       string
		inputSlice models.PostSlice
		inputCount int64
		expected   *graphql.PostsPayload
	}{
		{
			name: SuccessCase,
			inputSlice: []*models.Post{
				{
					ID:       1,
					AuthorID: 2,
					Content:  "content",
				},
				{
					ID:       2,
					AuthorID: 2,
					Content:  "content",
				},
			},
			inputCount: 2,
			expected: &graphql.PostsPayload{
				Posts: []*graphql.Post{
					{
						ID:       "1",
						AuthorID: "2",
						Content:  "content",
					},
					{
						ID:       "2",
						AuthorID: "2",
						Content:  "content",
					},
				},
				Total: 2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(
				t, test.expected, PostsToGraphqlPostsPayload(test.inputSlice, test.inputCount))
		})
	}
}
