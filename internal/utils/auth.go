package utils

import (
	"fmt"
	"os"
	"strings"
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

func ExtractTokenFromHeader(authorizationHeader string) string {
	if strings.HasPrefix(authorizationHeader, "Bearer ") {
		return strings.TrimPrefix(authorizationHeader, "Bearer ")
	}
	return ""
}

func ParseAndValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if err := ValidateSigningMethod(jwtToken); err != nil {
			return nil, err
		}
		return GetSecretKey()
	})
}

func ValidateSigningMethod(jwtToken *jwt.Token) error {
	if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
		return fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
	}
	return nil
}

func GetSecretKey() ([]byte, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("secret key not found in environment variables")
	}
	return []byte(secretKey), nil
}
