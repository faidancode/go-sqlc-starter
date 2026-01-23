package category

import "time"

// --- REQUEST DTO ---

type CreateCategoryRequest struct {
	// Tambahkan binding:"required" agar Gin menjalankan validasi otomatis
	Name        string `json:"name" binding:"required" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
	ImageUrl    string `json:"imageUrl" validate:"omitempty,url"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
	ImageUrl    string `json:"imageUrl" validate:"omitempty,url"`
	IsActive    *bool  `json:"is_active" binding:"required" validate:"required"`
}

type ListCategoryRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Limit   int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search  string `form:"search"`
	SortBy  string `form:"sortBy" binding:"omitempty,oneof=name created_at"`
	SortDir string `form:"sortDir" binding:"omitempty,oneof=asc desc"`
}

// --- RESPONSE DTO ---

type CategoryPublicResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	ImageUrl    string `json:"imageUrl,omitempty"`
}

type CategoryAdminResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	ImageUrl    string     `json:"imageUrl"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}
