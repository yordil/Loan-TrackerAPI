package infrastructure

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(id string , duration time.Duration , secret string ) (string, error) {


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
		"exp": time.Now().Add(time.Minute * duration).Unix(),
	})

	tokenString, err := token.SignedString([]byte("Secret"))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}


func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
