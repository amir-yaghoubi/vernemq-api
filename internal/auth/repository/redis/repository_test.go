package redis_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/amir-yaghoubi/vernemq-auth/internal/auth"
	redisRepo "github.com/amir-yaghoubi/vernemq-auth/internal/auth/repository/redis"
	"github.com/go-redis/redis/v7"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func testUser() auth.User {
	now := time.Now()

	var allowedRetain bool = true
	var notAllowedRetain bool = false

	var qos1 uint8 = 1
	var qos2 uint8 = 2

	return auth.User{
		Username:   "username",
		Password:   "super secret password ^^",
		Mountpoint: "github",
		PublishACL: []auth.PublishACL{
			{Pattern: "foo/bar/+", MaxQos: &qos1, AllowedRetain: &allowedRetain},
			{Pattern: "foo/baz/#", MaxQos: &qos2, AllowedRetain: &notAllowedRetain},
		},
		SubACL: []auth.SubACL{
			{Pattern: "foo/bar/+", MaxQos: &qos1},
			{Pattern: "foo/baz/#", MaxQos: &qos2},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func createRepo() (*redisRepo.Repository, func(), error) {
	s, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	repo := redisRepo.New(client)

	r := repo.(*redisRepo.Repository)
	return r, s.Close, nil
}

func TestSimpleUsage(t *testing.T) {
	repo, close, err := createRepo()
	if close != nil {
		defer close()
	}

	assert.NoError(t, err)

	user := testUser()
	err = repo.Set(&user)
	assert.NoError(t, err)

	u, err := repo.Get(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, cmp.Equal(&user, u), true)

	invalidU, err := repo.Get("invalid username")
	assert.EqualError(t, err, auth.ErrUserNotFound.Error())
	assert.Nil(t, invalidU)

	isDeleted, err := repo.Delete("username")
	assert.NoError(t, err)
	assert.Equal(t, true, isDeleted)

	isDeleted, err = repo.Delete("invalid username")
	assert.NoError(t, err)
	assert.Equal(t, false, isDeleted)
}
