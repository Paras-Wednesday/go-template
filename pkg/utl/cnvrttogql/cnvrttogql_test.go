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
			},
			expected: &graphql.Author{
				ID:        "29",
				FirstName: "First",
				LastName:  "",
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
