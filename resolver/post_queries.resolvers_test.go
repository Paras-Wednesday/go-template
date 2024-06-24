package resolver_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go-template/daos"
	fm "go-template/gqlmodels"
	"go-template/models"
	"go-template/resolver"
)

func TestPostByID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		wantResp *fm.Post
		wantErr  bool
		init     func() *gomonkey.Patches
	}{
		{
			name: "Should find the post",
			id:   "1",
			wantResp: &fm.Post{
				ID:       "1",
				AuthorID: "1",
				Content:  "content",
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(
					daos.FindPostForAuthorByID,
					func(ctx context.Context, authorID int, id int) (*models.Post, error) {
						return &models.Post{
							ID:       1,
							AuthorID: 1,
							Content:  "content",
						}, nil
					})
			},
		},
		{
			name:     "Should not find the post",
			id:       "hi",
			wantResp: nil,
			wantErr:  true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(
					daos.FindPostForAuthorByID,
					func(ctx context.Context, authorID int, id int) (*models.Post, error) {
						return nil, fmt.Errorf("can't find post")
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
			response, err := resolver.Query().PostByID(context.Background(), test.id)
			assert.Equal(t, test.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, test.wantResp, response)
			}
		})
	}
}

// nolint: funlen
func TestAllPostByAuthor(t *testing.T) {
	tests := []struct {
		name       string
		authorID   string
		pagination fm.Pagination
		wantResp   *fm.PostsPayload
		wantErr    bool
		init       func() *gomonkey.Patches
	}{
		{
			name:     "Should get 2 posts from author",
			authorID: "44",
			pagination: fm.Pagination{
				Limit: 5,
				Page:  2,
			},
			wantResp: &fm.PostsPayload{
				Posts: []*fm.Post{
					{
						ID:       "23",
						AuthorID: "44",
					},
					{
						ID:       "24",
						AuthorID: "44",
					},
				},
				Total: 10,
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(
					daos.FindAllPostBylAuthorWithCount,
					func(ctx context.Context,
						authorID int,
						queries ...qm.QueryMod,
					) (models.PostSlice, int64, error) {
						return []*models.Post{
							{
								ID:       23,
								AuthorID: 44,
							},
							{
								ID:       24,
								AuthorID: 44,
							},
						}, 10, nil
					})
			},
		},
		{
			name:     "Should return posts by given author",
			authorID: "hi",
			pagination: fm.Pagination{
				Limit: 2,
				Page:  3,
			},
			wantResp: &fm.PostsPayload{
				Posts: []*fm.Post{},
				Total: 0,
			},
			wantErr: true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(
					daos.FindAllPostBylAuthorWithCount,
					func(ctx context.Context,
						authorID int,
						queries ...qm.QueryMod,
					) (models.PostSlice, int64, error) {
						return models.PostSlice{}, 0, fmt.Errorf("no such author")
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
			response, err := resolver.Query().AllPostByAuthor(
				context.Background(), test.pagination)
			assert.Equal(t, test.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, test.wantResp, response)
			}
		})
	}
}
