package response

import (
	"net/http"

	"event-coming/internal/domain"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *PaginationMeta `json:"meta"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created sends a created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a no content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, details interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "validation_error",
			Message: "Validation failed",
			Details: details,
		},
	})
}

// HandleDomainError handles domain errors and sends appropriate responses
func HandleDomainError(c *gin.Context, err error) {
	switch err {
	case domain.ErrNotFound:
		Error(c, http.StatusNotFound, "not_found", "Resource not found")
	case domain.ErrUnauthorized:
		Error(c, http.StatusUnauthorized, "unauthorized", "Unauthorized")
	case domain.ErrForbidden:
		Error(c, http.StatusForbidden, "forbidden", "Forbidden")
	case domain.ErrConflict:
		Error(c, http.StatusConflict, "conflict", "Resource already exists")
	case domain.ErrInvalidInput:
		Error(c, http.StatusBadRequest, "invalid_input", "Invalid input")
	case domain.ErrInvalidCredentials:
		Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid credentials")
	case domain.ErrTokenExpired:
		Error(c, http.StatusUnauthorized, "token_expired", "Token expired")
	case domain.ErrInvalidToken:
		Error(c, http.StatusUnauthorized, "invalid_token", "Invalid token")
	default:
		Error(c, http.StatusInternalServerError, "internal_error", "Internal server error")
	}
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, data interface{}, page, perPage int, total int64) {
	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: &PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
