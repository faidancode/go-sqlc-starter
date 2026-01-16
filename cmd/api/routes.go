package main

import (
	"go-sqlc-starter/internal/auth"
	"go-sqlc-starter/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ControllerRegistry struct {
	Auth *auth.Controller
	// Category *category.Controller
	// Product  *product.Controller
	// Cart     *cart.Controller
	// Order    *order.Controller
}

func setupRoutes(r *gin.Engine, reg ControllerRegistry) {
	r.Use(middleware.RequestID())

	v1 := r.Group("/api/v1")
	{
		// Auth Routes (Public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", reg.Auth.Login)
			auth.POST("/logout", reg.Auth.Logout)
			auth.POST("/register", reg.Auth.Register)
		}

	}
}
