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
