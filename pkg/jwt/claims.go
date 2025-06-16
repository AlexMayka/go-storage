package jwt

import "github.com/golang-jwt/jwt/v5"

type AuthClaims struct {
	UserID    string
	RoleID    string
	CompanyID string
	jwt.RegisteredClaims
}
