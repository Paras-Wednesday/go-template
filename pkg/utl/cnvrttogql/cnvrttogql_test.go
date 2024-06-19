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
			assert.Equal(t, test.expected, PostToGraphqlPost(test.input))
		})
	}
}

func TestPostsToGraphqlPostsPayload(t *testing.T) {
	tests := []struct {
		name       string
		inputPosts models.PostSlice
		inputCount int64
		expected   *graphql.PostsPayload
	}{
		{
			name:       "Should return empty posts",
			inputPosts: []*models.Post{},
			inputCount: 5,
			expected: &graphql.PostsPayload{
				Posts: []*graphql.Post{},
				Total: 5,
			},
		},
		{
			name: "Should return two posts",
			inputPosts: []*models.Post{
				{
					ID:       1,
					AuthorID: 3,
					Content:  "first post",
				},
				{
					ID:       2,
					AuthorID: 3,
					Content:  "second post",
				},
			},
			inputCount: 2,
			expected: &graphql.PostsPayload{
				Posts: []*graphql.Post{
					{
						ID:       "1",
						AuthorID: "3",
						Content:  "first post",
					},
					{
						ID:       "2",
						AuthorID: "3",
						Content:  "second post",
					},
				},
				Total: 2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := PostsToGraphqlPostsPayload(test.inputPosts, test.inputCount)
			assert.Equal(t, test.expected, got)
		})
	}
}
