package types

import (
	"github.com/syronz/dict"

	"github.com/dgrijalva/jwt-go"
)

// JWTClaims for JWT
type JWTClaims struct {
	Username string    `json:"username"`
	ID       RowID     `json:"id"`
	Lang     dict.Lang `json:"language"`
	jwt.StandardClaims
}
