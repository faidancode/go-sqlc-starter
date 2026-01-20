package main

import (
	"go-sqlc-starter/internal/api/v1/address"
	"go-sqlc-starter/internal/api/v1/auth"
	"go-sqlc-starter/internal/api/v1/cart"
	"go-sqlc-starter/internal/api/v1/category"
	"go-sqlc-starter/internal/api/v1/product"
	"go-sqlc-starter/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ControllerRegistry struct {
	Auth     *auth.Controller
	Category *category.Controller
	Product  *product.Controller
	Cart     *cart.Controller
	Address  *address.Controller
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

		prod := v1.Group("/products")
		{
			prod.GET("", reg.Product.GetPublicList)
			prod.GET("/:id", reg.Product.GetByID)
		}

		adminProduct := v1.Group("/admin/products")
		adminProduct.Use(middleware.AuthMiddleware())
		adminProduct.Use(middleware.RoleMiddleware("ADMIN", "SUPERADMIN"))
		{
			adminProduct.GET("", reg.Product.GetAdminList)
			adminProduct.POST("", reg.Product.Create)
			adminProduct.PUT("/:id", reg.Product.Update)
			adminProduct.DELETE("/:id", reg.Product.Delete)
			adminProduct.PATCH("/:id/restore", reg.Product.Restore)
		}

		cart := v1.Group("/cart")
		cart.Use(middleware.AuthMiddleware())
		{
			cart.POST("", reg.Cart.Create)
			cart.GET("", reg.Cart.Detail)
			cart.GET("/count", reg.Cart.Count)
			cart.DELETE("", reg.Cart.Delete)
		}

		cartItems := v1.Group("/cart-items")
		cartItems.Use(middleware.AuthMiddleware())
		{
			cartItems.PUT("/:id", reg.Cart.UpdateQty)
			cartItems.POST("/:id/increment", reg.Cart.Increment)
			cartItems.POST("/:id/decrement", reg.Cart.Decrement)
			cartItems.DELETE("/:id", reg.Cart.DeleteItem)
		}

		address := v1.Group("/address")
		address.Use(middleware.AuthMiddleware())
		{
			address.GET("/:user_id", reg.Address.List)
			address.POST("", reg.Address.Create)
			address.PUT("/:id", reg.Address.Update)
			address.DELETE("", reg.Address.Delete)
		}

		adminAddress := v1.Group("/admin/address")
		adminAddress.Use(middleware.AuthMiddleware())
		adminAddress.Use(middleware.RoleMiddleware("ADMIN", "SUPERADMIN"))
		{
			adminAddress.GET("", reg.Address.ListAdmin)
		}

	}
}
