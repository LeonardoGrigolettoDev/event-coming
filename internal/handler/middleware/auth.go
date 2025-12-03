package middleware

import (
	"strings"

	"event-coming/internal/config"
	"event-coming/internal/domain"
	"event-coming/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 401, "unauthorized", "Missing authorization header")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, 401, "unauthorized", "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.AccessSecret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, 401, "unauthorized", "Invalid token")
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, 401, "unauthorized", "Invalid token claims")
			c.Abort()
			return
		}

		// Set user info in context
		if userIDStr, ok := claims["user_id"].(string); ok {
			if userID, err := uuid.Parse(userIDStr); err == nil {
				c.Set("user_id", userID)
			}
		}

		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		}

		if orgIDStr, ok := claims["organization_id"].(string); ok {
			if orgID, err := uuid.Parse(orgIDStr); err == nil {
				c.Set("organization_id", orgID)
			}
		}

		if role, ok := claims["role"].(string); ok {
			c.Set("role", domain.UserRole(role))
		}

		c.Next()
	}
}

// RequireOrgAccess checks if the user has access to the organization
func RequireOrgAccess(requiredRole domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Error(c, 403, "forbidden", "No role found")
			c.Abort()
			return
		}

		userRole := role.(domain.UserRole)

		// Super admin can access everything
		if userRole == domain.UserRoleSuperAdmin {
			c.Next()
			return
		}

		// Check role hierarchy
		if !hasPermission(userRole, requiredRole) {
			response.Error(c, 403, "forbidden", "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasPermission(userRole, requiredRole domain.UserRole) bool {
	roleHierarchy := map[domain.UserRole]int{
		domain.UserRoleSuperAdmin:  6,
		domain.UserRoleOrgOwner:    5,
		domain.UserRoleOrgAdmin:    4,
		domain.UserRoleOrgManager:  3,
		domain.UserRoleOrgOperator: 2,
		domain.UserRoleOrgViewer:   1,
	}

	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}
