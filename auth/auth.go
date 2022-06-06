package idxauth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	*jwt.StandardClaims
}

type Auth interface {
	ExtractTokenMetadata(token *jwt.Token) (*Token, error)
	CreateToken(secret string, day time.Duration) string
	Validate(tokenString string, secret string) (*jwt.Token, error)
}
