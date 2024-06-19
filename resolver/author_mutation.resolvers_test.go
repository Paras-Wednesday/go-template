package resolver_test

import (
	"context"
	"fmt"
	"go-template/daos"
	fm "go-template/gqlmodels"
	"go-template/models"
	"go-template/resolver"
	"go-template/testutls"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func TestCreateAuthor(t *testing.T) {
	tests := []struct {
		name     string
		req      fm.AuthorCreateInput
		wantResp *fm.Author
		wantErr  bool
		init     func() *gomonkey.Patches
	}{
		{
			name: "Should create Author",
			req: fm.AuthorCreateInput{
				FirstName: "First",
				LastName:  "Last",
			},
			wantResp: &fm.Author{
				ID:        fmt.Sprintf("%d", testutls.MockID),
				FirstName: "First",
				LastName:  "Last",
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.CreateAuthor, func(ctx context.Context, author models.Author) (models.Author, error) {
					return *testutls.MockAuthor(), nil
				})
			},
		},
		{
			name: "Should not create Author",
			req: fm.AuthorCreateInput{
				FirstName: "First",
			},
			wantResp: nil,
			wantErr:  true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.CreateAuthor, func(ctx context.Context, author models.Author) (models.Author, error) {
					return models.Author{}, fmt.Errorf("Last name is required")
				})
			},
		},
	}
	resolver := resolver.Resolver{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup and tear down monkey patch
			patch := test.init()
			defer func() {
				if patch != nil {
					patch.Reset()
				}
			}()
			// sleep to reset the monkey patch
			time.Sleep(10 * time.Millisecond)

			response, err := resolver.Mutation().CreateAuthor(context.Background(), test.req)
			assert.Equal(t, test.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, test.wantResp, response)
			}
		})
	}
}

// nolint: funlen
func TestUpdateAuthor(t *testing.T) {
	tests := []struct {
		name     string
		req      fm.AuthorUpdateInput
		wantResp *fm.Author
		wantErr  bool
		init     func() *gomonkey.Patches
	}{
		{
			name: "Should update Author",
			req: fm.AuthorUpdateInput{
				ID:        "1",
				FirstName: &testutls.MockAuthor().FirstName,
			},
			wantResp: &fm.Author{
				ID:        "1",
				FirstName: testutls.MockAuthor().FirstName,
				LastName:  testutls.MockAuthor().LastName.String,
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindAuthorByID, func(ctx context.Context, id int) (*models.Author, error) {
					return testutls.MockAuthor(), nil
				}).ApplyFunc(daos.UpdateAuthor, func(ctx context.Context, author models.Author) (models.Author, error) {
					return *testutls.MockAuthor(), nil
				})
			},
		},
		{
			name: "Should not update Author",
			req: fm.AuthorUpdateInput{
				ID:        "fjdk",
				FirstName: &testutls.MockAuthor().FirstName,
			},
			wantResp: nil,
			wantErr:  true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindAuthorByID, func(ctx context.Context, id int) (*models.Author, error) {
					return nil, fmt.Errorf("no such author")
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
			time.Sleep(10 * time.Millisecond)
			resp, err := resolver.Mutation().UpdateAuthor(context.Background(), test.req)
			assert.Equal(t, test.wantErr, err != nil, "expected: %t got: %v", test.wantErr, err)
			if err == nil {
				assert.Equal(t, test.wantResp, resp)
			}
		})
	}
}

// nolint: funlen
func TestDeleteAuthor(t *testing.T) {
	tests := []struct {
		name     string
		req      fm.AuthorDeleteInput
		wantResp *fm.Author
		wantErr  bool
		init     func() *gomonkey.Patches
	}{
		{
			name: "Should delete Author",
			req: fm.AuthorDeleteInput{
				ID: "33",
			},
			wantResp: &fm.Author{
				ID:        "33",
				FirstName: "fname",
				LastName:  "lname",
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindAuthorByID, func(ctx context.Context, id int) (*models.Author, error) {
					return &models.Author{
						ID:        33,
						FirstName: "fname",
						LastName:  null.StringFrom("lname"),
					}, nil
				}).ApplyFunc(daos.DeleteAuthor, func(ctx context.Context, author models.Author) (int64, error) {
					return 1, nil
				})
			},
		},
		{
			name: "Should not delete Author",
			req: fm.AuthorDeleteInput{
				ID: "fjdk",
			},
			wantResp: nil,
			wantErr:  true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindAuthorByID, func(ctx context.Context, id int) (*models.Author, error) {
					return nil, fmt.Errorf("no such author")
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
			time.Sleep(10 * time.Millisecond)
			resp, err := resolver.Mutation().DeleteAuthor(context.Background(), test.req)
			assert.Equal(t, test.wantErr, err != nil, "expected: %t got: %v", test.wantErr, err)
			if err == nil {
				assert.Equal(t, test.wantResp, resp)
			}
		})
	}
}
