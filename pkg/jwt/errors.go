package jwt

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token is expired")
	ErrEmptyUserID  = errors.New("empty user ID")
	ErrEmptyRoleID  = errors.New("empty role ID")
)
