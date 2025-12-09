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

// RoleHierarchy maps roles to their permission levels
var RoleHierarchy = map[domain.UserRole]int{
	domain.UserRoleSuperAdmin:    100,
	domain.UserRoleEntityOwner:   50,
	domain.UserRoleEntityAdmin:   40,
	domain.UserRoleEntityManager: 30,
	domain.UserRoleEntityViewer:  10,
}

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

		if entityIDStr, ok := claims["entity_id"].(string); ok {
			if orgID, err := uuid.Parse(entityIDStr); err == nil {
				c.Set("entity_id", orgID)
			}
		}

		if role, ok := claims["role"].(string); ok {
			c.Set("role", domain.UserRole(role))
		}

		c.Next()
	}
}

// RequireRole checks if the user has at least the required role level
func RequireRole(requiredRole domain.UserRole) gin.HandlerFunc {
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
		if !HasPermission(userRole, requiredRole) {
			response.Error(c, 403, "forbidden", "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// HasPermission checks if userRole has at least the permission level of requiredRole
func HasPermission(userRole, requiredRole domain.UserRole) bool {
	return RoleHierarchy[userRole] >= RoleHierarchy[requiredRole]
}

// RequireEntityAccess validates that the user has access to the entity in the route
// This should be used on routes that have :entity parameter
func RequireEntityAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get entity_id from JWT context
		entityID, exists := c.Get("entity_id")
		if !exists {
			response.Error(c, 403, "forbidden", "No entity access")
			c.Abort()
			return
		}

		// Get role from context
		role, roleExists := c.Get("role")

		// Super admin can access any entity
		if roleExists {
			userRole := role.(domain.UserRole)
			if userRole == domain.UserRoleSuperAdmin {
				c.Next()
				return
			}
		}

		// Check if route has :entity parameter and validate ownership
		entityParam := c.Param("entity")
		if entityParam != "" {
			routeEntityID, err := uuid.Parse(entityParam)
			if err != nil {
				response.Error(c, 400, "bad_request", "Invalid entity ID")
				c.Abort()
				return
			}

			userEntityID := entityID.(uuid.UUID)
			if routeEntityID != userEntityID {
				response.Error(c, 403, "forbidden", "Access denied to this entity")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireOwnerOrAdmin ensures user is owner or admin to perform sensitive operations
func RequireOwnerOrAdmin() gin.HandlerFunc {
	return RequireRole(domain.UserRoleEntityAdmin)
}
