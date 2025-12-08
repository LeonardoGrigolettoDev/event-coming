package handler

import (
	"net/http"

	"event-coming/internal/dto"
	"event-coming/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler contém as dependências para handlers de auth
type AuthHandler struct {
	authService service.AuthService // ← Dependência injetada
}

// NewAuthHandler cria um novo AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login processa POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Register processa POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	// 1. Parse + Validação do JSON
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	// 2. Chamar o service
	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		// TODO: tratar erros específicos (email duplicado, etc.)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 3. Retornar resposta de sucesso
	c.JSON(http.StatusCreated, result)
}

// Refresh processa POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := h.authService.Refresh(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired token",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Logout processa POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// ForgotPassword processa POST /auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := h.authService.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		// Mesmo em caso de erro, retornamos sucesso genérico por segurança
		c.JSON(http.StatusOK, gin.H{
			"message": "If an account with this email exists, a password reset link has been sent.",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ResetPassword processa POST /auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := h.authService.ResetPassword(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidToken {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid or expired reset token",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reset password",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
