package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidJWTToken = errors.New("invalid JWT token")
)

func CreateToken(secretKey, username string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		},
	)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func parseToken(tokenString string, secretKey []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}

func VerifyToken(tokenString, secretKey string) (string, error) {
	
	token, err := parseToken(tokenString, []byte(secretKey))

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrInvalidJWTToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidJWTToken
	}

	username := claims["username"].(string)
	return username, nil
}


