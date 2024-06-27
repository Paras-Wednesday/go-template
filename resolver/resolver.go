package resolver

import (
	"sync"

	fm "go-template/gqlmodels"
	pm "go-template/post-model"
)

// This file will
// not be
// regenerated
// automatically.
//
// It serves as
// dependency
// injection for
// your app, add any
// dependencies you
// require here.

//go:generate mockgen -source=resolver.go -destination=./mocks/mock.go PostDAO
type PostDAO interface {
	CreatePost(post pm.PostModel) (pm.PostModel, error)
	GetPost(id int) (pm.PostModel, error)
}

// Resolver ...
type Resolver struct {
	sync.Mutex
	Observers map[string]chan *fm.User
	PostDAO   PostDAO
}
