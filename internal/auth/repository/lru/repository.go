package cache

import (
	"github.com/amir-yaghoubi/vernemq-auth/internal/auth"
	lru "github.com/hashicorp/golang-lru"
)

// New Create a new Repository with lru cache
func New(repo auth.Repository, cacheSize int) (auth.Repository, error) {
	cache, err := lru.NewARC(cacheSize)
	if err != nil {
		return nil, err
	}

	return &Repository{repository: repo, cache: cache}, nil
}

// Repository wrapped with lru cache
type Repository struct {
	repository auth.Repository
	cache      *lru.ARCCache
}

func (r *Repository) Get(username string) (*auth.User, error) {
	cacheUser, ok := r.cache.Get(username)
	if ok {
		user := cacheUser.(auth.User)
		return &user, nil
	}

	user, err := r.repository.Get(username)
	if err != nil {
		return nil, err
	}
	r.cache.Add(username, *user)

	return user, nil
}

func (r *Repository) Set(user *auth.User) error {
	err := r.repository.Set(user)
	if err != nil {
		return err
	}

	r.cache.Add(user.Username, *user)
	return nil
}

func (r *Repository) Delete(username string) (bool, error) {
	r.cache.Remove(username)
	return r.repository.Delete(username)
}
