package auth

import (
	"go-sqlc-starter/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service // Perbaikan: Gunakan Interface, bukan pointer ke Interface
}

func NewController(s Service) *Controller {
	return &Controller{service: s}
}

func (ctrl *Controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Input tidak valid", err.Error())
		return
	}

	token, userResp, err := ctrl.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "AUTH_FAILED", "Email atau password salah", nil)
		return
	}

	// Cookie configuration
	c.SetCookie("access_token", token, 86400, "/", "", false, true)

	response.Success(c, http.StatusOK, userResp, nil)
}

func (ctrl *Controller) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Input tidak valid", err.Error())
		return
	}

	res, err := ctrl.service.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "REGISTER_FAILED", err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, res, nil)
}

func (ctrl *Controller) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	response.Success(c, http.StatusOK, "Logout berhasil", nil)
}
