package auth

import "errors"

// ErrUserNotFound occure when user does not exists in our database
var ErrUserNotFound = errors.New("user not found")

// ErrUnAuthorizeAccess occure when user does not have access for that operation
var ErrUnAuthorizeAccess = errors.New("unauthoirzed access")
