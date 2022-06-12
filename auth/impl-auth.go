package idxauth

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type authService struct{}

func (*authService) CreateToken(secret string, day time.Duration) string {

	expiresAt := time.Now().Add(time.Hour * 24 * day).Unix()
	tk := &Token{

		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println(err)
	}

	return tokenString

}

func (a *authService) Validate(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	log.Println(err)

	if err != nil {
		return nil, err
	} else {
		return token, nil
	}

}

func (*authService) ExtractTokenMetadata(token *jwt.Token) (*Token, error) {
	_, ok := token.Claims.(jwt.MapClaims)
	var dataToken Token

	if ok && token.Valid {

		return &dataToken, nil
	} else {
		return nil, fmt.Errorf("")
	}

}

func ImplAuthService() Auth {
	return &authService{}
}
