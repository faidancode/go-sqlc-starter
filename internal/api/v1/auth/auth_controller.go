package auth

import (
	"go-sqlc-starter/internal/pkg/platform"
	"go-sqlc-starter/internal/pkg/response"
	"net/http"
	"os"

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
		// Response Error Seragam
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Input tidak valid", err.Error())
		return
	}

	clientHeader := c.GetHeader("X-Client-Type")
	userAgent := c.GetHeader("User-Agent")
	clientType := platform.ResolveClientType(clientHeader, userAgent)

	token, refreshToken, userResp, err := ctrl.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// Response Error Seragam
		response.Error(c, http.StatusUnauthorized, "AUTH_FAILED", "Email atau password salah", nil)
		return
	}
	isProd := os.Getenv("APP_ENV") == "production"

	if platform.IsWebClient(clientType) {
		c.SetCookie(
			"access_token",
			token,
			86400,
			"/",
			"",
			isProd,
			true,
		)

		c.SetCookie(
			"refresh_token",
			refreshToken,
			3600*24*7,
			"/",
			"",
			isProd,
			true)
	}

	responseData := gin.H{
		"user":          userResp,
		"access_token":  token,
		"refresh_token": refreshToken,
	}

	response.Success(c, http.StatusOK, responseData, nil)
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
