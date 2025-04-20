package jwtpkg

import (
	"os"
	"time"

	"github.com/Gergenus/AuthService/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()

	tkn, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tkn, nil
}
