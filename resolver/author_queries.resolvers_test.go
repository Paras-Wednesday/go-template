package resolver_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go-template/daos"
	fm "go-template/gqlmodels"
	"go-template/models"
	"go-template/resolver"
)

// nolint: funlen
func TestAuthor(t *testing.T) {
	now := time.Now()
	nowMilli := int(time.Now().UnixMilli())
	tests := []struct {
		name       string
		inputID    string
		wantAuthor *fm.Author
		wantErr    bool
		init       func() *gomonkey.Patches
	}{
		{
			name:    "Should get author with id 2",
			inputID: "2",
			wantAuthor: &fm.Author{
				ID:        "2",
				FirstName: "John",
				LastName:  "Doe",
				CreatedAt: &nowMilli,
				UpdatedAt: &nowMilli,
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindAuthorByID, func(ctx context.Context, id int) (*models.Author, error) {
					return &models.Author{
						ID:        2,
						FirstName: "John",
						LastName:  null.StringFrom("Doe"),
						CreatedAt: null.TimeFrom(now),
						UpdatedAt: null.TimeFrom(now),
					}, nil
				})
			},
		},
		{
			name:       "Should return error for input id 'hi'",
			inputID:    "hi",
			wantAuthor: nil,
			wantErr:    true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindAuthorByID, func(ctx context.Context, id int) (*models.Author, error) {
					return nil, fmt.Errorf("no such author for id 'hi'")
				})
			},
		},
	}

	resolver := resolver.Resolver{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			patch := test.init()
			defer func() {
				if patch != nil {
					patch.Reset()
				}
			}()

			author, err := resolver.Query().Author(context.Background(), test.inputID)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err,
			)
			assert.Equal(t, author, test.wantAuthor)
		})
	}
}

// nolint: funlen
func TestAuthors(t *testing.T) {
	tests := []struct {
		name        string
		input       fm.Pagination
		wantPayload *fm.AuthorsPayload
		wantErr     bool
		init        func() *gomonkey.Patches
	}{
		{
			name: "Should return 2 authors out of 20 ",
			input: fm.Pagination{
				Limit: 2,
				Page:  1,
			},
			wantPayload: &fm.AuthorsPayload{
				Authors: []*fm.Author{
					{
						ID:        "1",
						FirstName: "author_first01",
						LastName:  "author_last01",
					},
					{
						ID:        "2",
						FirstName: "author_first02",
						LastName:  "author_last02",
					},
				},
				Total: 20,
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(
					daos.GetAllAuthorsWithCount,
					func(ctx context.Context, queries ...qm.QueryMod) (models.AuthorSlice, int64, error) {
						return []*models.Author{
							{
								ID:        1,
								FirstName: "author_first01",
								LastName:  null.StringFrom("author_last01"),
							},
							{
								ID:        2,
								FirstName: "author_first02",
								LastName:  null.StringFrom("author_last02"),
							},
						}, 20, nil
					})
			},
		},
		{
			name: "Should return sql error",
			input: fm.Pagination{
				Limit: 2,
				Page:  1,
			},
			wantPayload: nil,
			wantErr:     true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(
					daos.GetAllAuthorsWithCount,
					func(ctx context.Context, queries ...qm.QueryMod) (models.AuthorSlice, int64, error) {
						return models.AuthorSlice{}, 0, fmt.Errorf("sql error")
					})
			},
		},
		{
			name: "Should return validation error for negative limit",
			input: fm.Pagination{
				Limit: -2,
				Page:  1,
			},
			wantPayload: nil,
			wantErr:     true,
			init: func() *gomonkey.Patches {
				return nil
			},
		},
		{
			name: "Should return validation error for negative page",
			input: fm.Pagination{
				Limit: 2,
				Page:  -1,
			},
			wantPayload: nil,
			wantErr:     true,
			init: func() *gomonkey.Patches {
				return nil
			},
		},
	}

	resolver := resolver.Resolver{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			patch := test.init()
			defer func() {
				if patch != nil {
					patch.Reset()
				}
			}()
			authorPayload, err := resolver.Query().Authors(context.Background(), test.input)
			assert.Equal(t, test.wantErr, err != nil,
				"wantErr: %t, got: %v", test.wantErr, err,
			)
			assert.Equal(t, authorPayload, test.wantPayload)
		})
	}
}
