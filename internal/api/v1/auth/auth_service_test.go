package auth_test

import (
	"context"
	"errors"
	"go-sqlc-starter/internal/api/v1/auth"
	authMock "go-sqlc-starter/internal/api/v1/mock/auth"
	"go-sqlc-starter/internal/dbgen"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := authMock.NewMockRepository(ctrl)
	service := auth.NewService(mockRepo)
	ctx := context.Background()

	pw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	t.Run("Success Login", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByEmail(ctx, "admin").
			Return(dbgen.GetUserByEmailRow{Email: "admin", Password: string(pw)}, nil)

		token, refreshToken, resp, err := service.Login(ctx, "admin", "password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, refreshToken)
		assert.Equal(t, "admin", resp.Email)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByEmail(ctx, "admin").
			Return(dbgen.GetUserByEmailRow{Email: "admin", Password: string(pw)}, nil)

		_, _, _, err := service.Login(ctx, "admin", "wrongpass")
		assert.Error(t, err)
	})
}

func TestService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := authMock.NewMockRepository(ctrl)
	service := auth.NewService(mockRepo)
	ctx := context.Background()

	t.Run("Success Register", func(t *testing.T) {
		req := auth.RegisterRequest{
			Email:    "user@example.com",
			Password: "password123",
		}

		mockRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(dbgen.CreateUserRow{
				Email: req.Email,
				Role:  "CUSTOMER",
			}, nil)

		resp, err := service.Register(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Email, resp.Email)
		assert.Equal(t, "CUSTOMER", resp.Role)
	})

	t.Run("Error Register", func(t *testing.T) {
		req := auth.RegisterRequest{
			Email:    "user@example.com",
			Password: "password123",
		}

		mockRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(dbgen.CreateUserRow{}, errors.New("duplicate email"))

		_, err := service.Register(ctx, req)
		assert.Error(t, err)
	})
}
