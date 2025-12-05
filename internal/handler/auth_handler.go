package handler

import (
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
	// 1. Parse do JSON → DTO
	// 2. Validação
	// 3. Chamar h.authService.Login(...)
	// 4. Retornar Response
}

// Register processa POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	// ...
}

// Refresh processa POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	// ...
}
