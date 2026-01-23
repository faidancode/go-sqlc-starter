package product

import "time"

type ListPublicRequest struct {
	Page       int
	Limit      int
	Search     string
	CategoryID string
	MinPrice   float64
	MaxPrice   float64
	SortBy     string
}

type ProductPublicResponse struct {
	ID           string  `json:"id"`
	CategoryName string  `json:"categoryName"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Price        float64 `json:"price"`
}

type ProductAdminResponse struct {
	ID           string    `json:"id"`
	CategoryName string    `json:"categoryName"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Price        float64   `json:"price"`
	Stock        int32     `json:"stock"`
	SKU          string    `json:"sku"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
}

type CreateProductRequest struct {
	CategoryID  string  `json:"categoryId" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int32   `json:"stock" binding:"required"`
	SKU         string  `json:"sku"`
	ImageUrl    string  `json:"imageUrl"`
}

type UpdateProductRequest struct {
	CategoryID  string  `json:"categoryId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
	SKU         string  `json:"sku"`
	ImageUrl    string  `json:"imageUrl"`
	IsActive    *bool   `json:"isActive"` // Gunakan pointer agar bisa membedakan false (bool) dan nil (tidak dikirim)
}
