package services

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"time"

	"github.com/Realwale/scribana/internal/config"
	"github.com/Realwale/scribana/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{Db: db}
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(user.ID), 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWTExpiresIn)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *AuthService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *AuthService) Authenticate(email, password string) (*models.User, error) {
	var user models.User
	if err := s.Db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := s.VerifyPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
