package helpers

import (
	"os"
	"p2-mini-project/src/entity"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func CheckHashPassword(hashedPass string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password)); err != nil {
		return err
	}
	return nil
}

func CreateJWT(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"fullname": user.Fullname,
		"user_id":  user.ID,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret_token := []byte(os.Getenv("JWT"))

	tokenString, err := token.SignedString(secret_token)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
