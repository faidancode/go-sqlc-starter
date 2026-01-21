package main

import (
	"go-sqlc-starter/internal/api/v1/address"
	"go-sqlc-starter/internal/api/v1/auth"
	"go-sqlc-starter/internal/api/v1/cart"
	"go-sqlc-starter/internal/api/v1/category"
	"go-sqlc-starter/internal/api/v1/order"
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
	Order    *order.Controller
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

		categories := v1.Group("/categories")
		{
			categories.GET("", reg.Category.ListPublic)
			categories.GET("/:id", reg.Category.GetByID)
		}

		adminCategories := categories.Group("/admin/categories")
		adminCategories.Use(middleware.AuthMiddleware())
		adminCategories.Use(middleware.RoleMiddleware("ADMIN", "SUPERADMIN"))
		{
			adminCategories.GET("", reg.Category.ListAdmin)
			adminCategories.POST("", reg.Category.Create)
			adminCategories.PUT("/:id", reg.Category.Update)
			adminCategories.DELETE("/:id", reg.Category.Delete)
			adminCategories.PATCH("/:id/restore", reg.Category.Restore)
		}

		products := v1.Group("/products")
		{
			products.GET("", reg.Product.GetPublicList)
			products.GET("/:id", reg.Product.GetByID)
		}

		adminProducts := v1.Group("/admin/products")
		adminProducts.Use(middleware.AuthMiddleware())
		adminProducts.Use(middleware.RoleMiddleware("ADMIN", "SUPERADMIN"))
		{
			adminProducts.GET("", reg.Product.GetAdminList)
			adminProducts.POST("", reg.Product.Create)
			adminProducts.PUT("/:id", reg.Product.Update)
			adminProducts.DELETE("/:id", reg.Product.Delete)
			adminProducts.PATCH("/:id/restore", reg.Product.Restore)
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

		// ========================
		// ORDER
		// ========================
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthMiddleware()) // Semua route order butuh login
		{
			// Customer Routes
			orders.POST("/checkout", reg.Order.Checkout)
			orders.GET("", reg.Order.List)
			orders.GET("/:id", reg.Order.Detail)
			orders.PATCH("/:id/cancel", reg.Order.Cancel)
			orders.PATCH("/:id/status", reg.Order.UpdateStatusByCustomer)

			// Admin Routes (Management)
			adminOrders := orders.Group("/admin")
			adminOrders.Use(middleware.RoleMiddleware("ADMIN", "SUPERADMIN"))
			{
				adminOrders.GET("", reg.Order.ListAdmin)
				adminOrders.PATCH("/:id/status", reg.Order.UpdateStatusByAdmin)
			}
		}

	}
}
