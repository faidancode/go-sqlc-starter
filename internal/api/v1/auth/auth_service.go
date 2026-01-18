package auth

import (
	"context"
	"database/sql"
	"fmt"
	"go-sqlc-starter/internal/dbgen"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=auth_service.go -destination=../mock/auth/auth_service_mock.go -package=mock
type Service interface {
	Register(ctx context.Context, req RegisterRequest) (AuthResponse, error)
	Login(ctx context.Context, email, password string) (string, AuthResponse, error)
	GetProfile(ctx context.Context, userID string) (AuthResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Login(ctx context.Context, email, password string) (string, AuthResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", AuthResponse{}, fmt.Errorf("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", AuthResponse{}, fmt.Errorf("invalid email or password")
	}

	tokenString, err := s.generateToken(user.ID.String(), user.Role.String)
	if err != nil {
		return "", AuthResponse{}, fmt.Errorf("failed to generate token")
	}

	return tokenString, AuthResponse{
		Email: user.Email,
		Role:  user.Role.String,
	}, nil
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (AuthResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to hash password")
	}

	user, err := s.repo.Create(ctx, dbgen.CreateUserParams{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  string(hashed),
		Role: sql.NullString{
			String: "customer",
			Valid:  true,
		},
	})
	if err != nil {
		return AuthResponse{}, fmt.Errorf("email already registered")
	}

	return AuthResponse{
		Email: user.Email,
		Role:  user.Role.String,
	}, nil
}

func (s *service) generateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *service) GetProfile(ctx context.Context, userID string) (AuthResponse, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("invalid user id format")
	}

	user, err := s.repo.GetByID(ctx, parsedID)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("user not found")
	}

	return AuthResponse{
		Email:     user.Email,
		Role:      user.Role.String,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
