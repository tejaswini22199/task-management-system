package services

import (
	"github.com/tejaswini22199/task-management-system/authservice/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const JWTSecret = "your-secret-key"

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterUserService(input RegisterInput) (map[string]interface{}, error) {
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	userID, err := repository.InsertUser(input.Name, input.Email, hashedPassword)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	token, err := generateToken(userID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return map[string]interface{}{
		"user_id": userID,
		"token":   token,
		"message": "User registered successfully",
	}, nil
}

func LoginUserService(input LoginInput) (map[string]interface{}, error) {
	userID, storedPassword, err := repository.GetUserByEmail(input.Email)
	if err != nil || !checkPasswordHash(input.Password, storedPassword) {
		return nil, errors.New("invalid credentials")
	}

	token, err := generateToken(userID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return map[string]interface{}{
		"user_id": userID,
		"token":   token,
		"message": "Login successful",
	}, nil
}
