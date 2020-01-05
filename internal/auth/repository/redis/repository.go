package redis

import (
	"errors"
	"github.com/amir-yaghoubi/vernemq-api/internal/auth"
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"
)

// New Create a new Redis Repository
func New(redisClient *redis.Client) auth.Repository {
	return &Repository{redis: redisClient}
}

// Repository redis implementation of auth repository
type Repository struct {
	redis *redis.Client
}

func (repo *Repository) key(username string) string {
	return "user:" + username
}

func (repo *Repository) encode(user *auth.User) ([]byte, error) {
	encoded, err := msgpack.Marshal(user)
	return encoded, err
}

func (repo *Repository) decode(bytes []byte) (*auth.User, error) {
	user := auth.User{}
	err := msgpack.Unmarshal(bytes, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *Repository) Get(username string) (*auth.User, error) {
	bytes, err := repo.redis.Get(repo.key(username)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}

	return repo.decode(bytes)
}

func (repo *Repository) Set(user *auth.User) error {
	encoded, err := repo.encode(user)
	if err != nil {
		return err
	}

	return repo.redis.Set(repo.key(user.Username), encoded, -1).Err()
}

func (repo *Repository) Delete(username string) (bool, error) {
	result, err := repo.redis.Del(repo.key(username)).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
