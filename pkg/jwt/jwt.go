package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenResponse struct {
	Token     string
	ExpiresAt time.Time
}

func CreateToken(userID, roleID, companyID string, secret []byte) (string, error) {
	tokenResp, err := CreateTokenWithExpiry(userID, roleID, companyID, secret)
	if err != nil {
		return "", err
	}
	return tokenResp.Token, nil
}

func CreateTokenWithExpiry(userID, roleID, companyID string, secret []byte) (*TokenResponse, error) {
	if userID == "" {
		return nil, ErrEmptyUserID
	}

	if roleID == "" {
		return nil, ErrEmptyRoleID
	}

	// super_admin может не иметь company_id
	// if companyID == "" {
	//	return nil, ErrEmptyCompanyID
	// }

	expiresAt := time.Now().Add(time.Hour * 72)
	claims := AuthClaims{
		RoleID:    roleID,
		UserID:    userID,
		CompanyID: companyID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}, nil
}

func ParseToken(tokenStr string, secret []byte) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, ErrTokenExpired
	}

	return claims, err
}
