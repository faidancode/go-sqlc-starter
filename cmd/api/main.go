package main

import (
	"database/sql"
	"go-sqlc-starter/internal/api/v1/auth"
	"go-sqlc-starter/internal/api/v1/category"
	"go-sqlc-starter/internal/api/v1/product"
	"go-sqlc-starter/internal/dbgen"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Load Environment
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 2. Database Connection
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer db.Close()

	// 3. Initialize SQLC Queries
	queries := dbgen.New(db)

	// 4. Initialize Modules (Dependency Injection)

	authRepo := auth.NewRepository(queries)
	authService := auth.NewService(authRepo)
	authController := auth.NewController(authService)

	categoryRepo := category.NewRepository(queries)
	categoryService := category.NewService(categoryRepo)
	categoryController := category.NewController(categoryService)

	productRepo := product.NewRepository(queries)
	productService := product.NewService(productRepo, categoryRepo)
	productController := product.NewController(productService)

	registry := ControllerRegistry{
		Auth:     authController,
		Category: categoryController,
		Product:  productController,
	}

	// 4. Jalankan Router
	r := gin.Default()
	setupRoutes(r, registry)

	// 7. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	r.Run(":" + port)
}
