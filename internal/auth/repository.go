package auth

type Repository interface {
	Get(username string) (*User, error)
	Set(user *User) error
	Delete(username string) (bool, error)
}
