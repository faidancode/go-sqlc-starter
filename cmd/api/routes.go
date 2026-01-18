package main

import (
	"go-sqlc-starter/internal/api/v1/auth"
	"go-sqlc-starter/internal/api/v1/category"
	"go-sqlc-starter/internal/api/v1/product"
	"go-sqlc-starter/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ControllerRegistry struct {
	Auth     *auth.Controller
	Category *category.Controller
	Product  *product.Controller
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

		cat := v1.Group("/categories")
		{
			cat.GET("", reg.Category.GetAll)
			cat.GET("/:id", reg.Category.GetByID)

			// Protected Admin Routes
			adminCat := cat.Group("")
			adminCat.Use(middleware.AuthMiddleware())
			adminCat.Use(middleware.RoleMiddleware("ADMIN", "SUPERADMIN"))
			{
				adminCat.POST("", reg.Category.Create)
				adminCat.PUT("/:id", reg.Category.Update)
				adminCat.DELETE("/:id", reg.Category.Delete)
				adminCat.PATCH("/:id/restore", reg.Category.Restore)
			}
		}

		// Product Routes
		prod := v1.Group("/products")
		{
			prod.GET("", reg.Product.GetPublicList)
			prod.GET("/:id", reg.Product.GetByID)

			// Protected Admin Routes
			adminProd := prod.Group("/admin")
			adminProd.Use(middleware.AuthMiddleware())
			adminProd.Use(middleware.RoleMiddleware("ADMIN", "SUPERADMIN"))
			{
				adminProd.GET("", reg.Product.GetAdminList)
				adminProd.POST("", reg.Product.Create)
				adminProd.PUT("/:id", reg.Product.Update)
				adminProd.DELETE("/:id", reg.Product.Delete)
				adminProd.PATCH("/:id/restore", reg.Product.Restore)
			}
		}

	}
}
