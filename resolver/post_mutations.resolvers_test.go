package resolver_test

import (
	"context"
	"fmt"
	"testing"

	gomock "go.uber.org/mock/gomock"

	model "go-template/post-model"
	"go-template/resolver"
	mock_resolver "go-template/resolver/mocks"
)

func TestCreatePostSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDAO := mock_resolver.NewMockPostDAO(ctrl)
	post := model.PostModel{
		ID:      1,
		Content: "content",
	}
	mockDAO.EXPECT().CreatePost(post).Return(post, nil)

	resolver := resolver.Resolver{
		PostDAO: mockDAO,
	}

	response, err := resolver.Mutation().CreatePost(context.Background(), post.ID, post.Content)
	if err != nil {
		t.Errorf("did not expect error, got: %v", err)
	}
	t.Logf("response: %+v", response)
}

func TestCreatePostFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDAO := mock_resolver.NewMockPostDAO(ctrl)
	post := model.PostModel{
		ID:      1,
		Content: "content",
	}
	mockDAO.EXPECT().CreatePost(post).Return(model.PostModel{}, fmt.Errorf("Error"))

	resolver := resolver.Resolver{
		PostDAO: mockDAO,
	}

	response, err := resolver.Mutation().CreatePost(context.Background(), post.ID, post.Content)
	if err == nil {
		t.Error("expected error got no error")
	}
	t.Logf("response: %+v, err: %v", response, err)
}
