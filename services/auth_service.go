package services

import (
	"errors"
	"os"
	"time"

	"smart-choice/models"
	"smart-choice/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	return repository.CreateUser(&user)
}

func Login(email, password string) (string, bool, error) {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return "", false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", false, errors.New("invalid credentials")
	}

	if user.TwoFA {
		return "", true, nil
	}

	token, err := GenerateJWT(&user)
	if err != nil {
		return "", false, err
	}

	return token, false, nil
}

func GenerateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
