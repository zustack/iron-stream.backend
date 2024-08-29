package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
