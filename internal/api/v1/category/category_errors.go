package category

import (
	"go-sqlc-starter/internal/pkg/apperror"
	"net/http"
)

var (
	ErrCategoryNotFound = apperror.New(
		apperror.CodeNotFound,
		"Category not found",
		http.StatusNotFound,
	)

	ErrInvalidCategoryID = apperror.New(
		apperror.CodeInvalidInput,
		"Invalid category ID format",
		http.StatusBadRequest,
	)
)
