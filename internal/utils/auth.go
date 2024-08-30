package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func MakeJWT(userId int64) (string, int64, error) {
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	expDuration := time.Hour * 24 * 30
	exp := now.Add(expDuration).Unix()
	claims["sub"] = userId
	claims["exp"] = exp
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	tokenString, err := tokenByte.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", 0, err
	}
	return tokenString, exp, nil
}
