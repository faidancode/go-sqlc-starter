package category

import (
	"go-sqlc-starter/internal/pkg/apperror"
	"go-sqlc-starter/internal/pkg/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(s Service) *Controller {
	return &Controller{service: s}
}

func (ctrl *Controller) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, total, err := ctrl.service.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		log.Printf("[Category GetAll Error]: %v", err)
		handleError(c, err)
		return
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	response.Success(c, http.StatusOK, data, &response.PaginationMeta{
		Total:      total,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   limit,
	})
}

func (ctrl *Controller) GetByID(c *gin.Context) {
	res, err := ctrl.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, res, nil)
}

func (ctrl *Controller) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		vErr := apperror.New(apperror.CodeInvalidInput, "Invalid input", http.StatusBadRequest)
		handleError(c, vErr)
		return
	}

	res, err := ctrl.service.Create(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, res, nil)
}

func (ctrl *Controller) Update(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperror.New(apperror.CodeInvalidInput, "Invalid input", http.StatusBadRequest))
		return
	}

	res, err := ctrl.service.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, res, nil)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	if err := ctrl.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, nil)
}

func (ctrl *Controller) Restore(c *gin.Context) {
	res, err := ctrl.service.Restore(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, res, nil)
}

// Helper function untuk mapping AppError ke Response
func handleError(c *gin.Context, err error) {
	httpErr := apperror.ToHTTP(err)
	response.Error(
		c,
		httpErr.Status,
		httpErr.Code,
		httpErr.Message,
		nil,
	)
	log.Println(err)
}
