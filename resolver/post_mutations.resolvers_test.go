package resolver_test

import (
	"context"
	"fmt"
	"go-template/daos"
	fm "go-template/gqlmodels"
	"go-template/models"
	"go-template/resolver"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name     string
		req      fm.PostCreateInput
		wantResp *fm.Post
		wantErr  bool
		init     func() *gomonkey.Patches
	}{
		{
			name: "should create post",
			req: fm.PostCreateInput{
				AuthorID: "1",
				Content:  "This is a good post",
			},
			wantResp: &fm.Post{
				ID:       "1",
				AuthorID: "1",
				Content:  "This is a good post",
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.CreatePost, func(ctx context.Context, post models.Post) (models.Post, error) {
					return models.Post{
						ID:       1,
						AuthorID: 1,
						Content:  "This is a good post",
					}, nil
				})
			},
		},
		{
			name: "should not create post",
			req: fm.PostCreateInput{
				Content: "This is a good post",
			},
			wantResp: nil,
			wantErr:  true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.CreatePost, func(ctx context.Context, post models.Post) (models.Post, error) {
					return models.Post{}, fmt.Errorf("author id is required")
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
			response, err := resolver.Mutation().CreatePost(context.Background(), test.req)
			assert.Equal(t, test.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, test.wantResp, response)
			}
		})
	}
}

// nolint: funlen
func TestUpdatePost(t *testing.T) {
	tests := []struct {
		name     string
		req      fm.PostUpdateInput
		wantResp *fm.Post
		wantErr  bool
		init     func() *gomonkey.Patches
	}{
		{
			name: "Should update the post",
			req: fm.PostUpdateInput{
				ID:      "1",
				Content: "this is updated content",
			},
			wantResp: &fm.Post{
				ID:       "1",
				AuthorID: "0",
				Content:  "this is updated content",
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindPostByID, func(ctx context.Context, id int) (*models.Post, error) {
					return &models.Post{
						ID:      1,
						Content: "This is original content",
					}, nil
				}).ApplyFunc(daos.UpdatePost, func(ctx context.Context, post models.Post) (models.Post, error) {
					return models.Post{
						ID:      1,
						Content: "this is updated content",
					}, nil
				})
			},
		},
		{
			name: "Should not update the post",
			req: fm.PostUpdateInput{
				ID:      "1",
				Content: "this is updated content",
			},
			wantResp: nil,
			wantErr:  true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindPostByID, func(ctx context.Context, id int) (*models.Post, error) {
					return nil, fmt.Errorf("can't find user")
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
			response, err := resolver.Mutation().UpdatePost(context.Background(), test.req)
			assert.Equal(t, test.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, test.wantResp, response)
			}
		})
	}
}

// nolint: funlen
func TestDeletePost(t *testing.T) {
	tests := []struct {
		name     string
		req      fm.PostDeleteInput
		wantResp *fm.Post
		wantErr  bool
		init     func() *gomonkey.Patches
	}{
		{
			name: "Should delete the post",
			req: fm.PostDeleteInput{
				ID: "1",
			},
			wantResp: &fm.Post{
				ID:       "23",
				AuthorID: "3",
				Content:  "content",
			},
			wantErr: false,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindPostByID, func(ctx context.Context, id int) (*models.Post, error) {
					return &models.Post{
						ID:       23,
						AuthorID: 3,
						Content:  "content",
					}, nil
				}).ApplyFunc(daos.DeletePost, func(ctx context.Context, post models.Post) (int64, error) {
					return 1, nil
				})
			},
		},
		{
			name: "Should not delete the post",
			req: fm.PostDeleteInput{
				ID: "1",
			},
			wantResp: nil,
			wantErr:  true,
			init: func() *gomonkey.Patches {
				return gomonkey.ApplyFunc(daos.FindPostByID, func(ctx context.Context, id int) (*models.Post, error) {
					return nil, fmt.Errorf("can't find user")
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
			response, err := resolver.Mutation().DeletePost(context.Background(), test.req)
			assert.Equal(t, test.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, test.wantResp, response)
			}
		})
	}
}
