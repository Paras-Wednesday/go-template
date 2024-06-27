package model

import (
	"errors"
	"time"
)

var (
	ErrPostExists = errors.New("post already exists")
	ErrNoPost     = errors.New("no such post")
)

type PostModel struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type PostRepo struct {
	m map[int]PostModel
}

func NewPostRepo() *PostRepo {
	return &PostRepo{
		m: map[int]PostModel{},
	}
}

func (r *PostRepo) CreatePost(post PostModel) (PostModel, error) {
	// check if it already exists
	if _, ok := r.m[post.ID]; ok {
		return PostModel{}, ErrPostExists
	}
	post.CreatedAt = time.Now()
	r.m[post.ID] = post
	return post, nil
}

func (r *PostRepo) GetPost(id int) (PostModel, error) {
	if post, ok := r.m[id]; !ok {
		return PostModel{}, ErrNoPost
	} else {
		return post, nil
	}
}
