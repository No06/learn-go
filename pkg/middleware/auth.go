package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Context keys.
const (
	ContextAccountID = "accountID"
	ContextRole      = "role"
)

// AuthConfig holds secret used to verify JWT.
type AuthConfig struct {
	Secret       string
	AllowedRoles []string
}

// JWTAuth returns a middleware that verifies JWT tokens and optionally enforces roles.
func JWTAuth(cfg AuthConfig) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(cfg.AllowedRoles))
	for _, role := range cfg.AllowedRoles {
		roleSet[role] = struct{}{}
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"message": "missing authorization"}})
			return
		}

		tokenString := strings.TrimSpace(authHeader[7:])

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(cfg.Secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"message": "invalid token"}})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"message": "invalid token claims"}})
			return
		}

		accountID, _ := claims["sub"].(string)
		role, _ := claims["role"].(string)

		if accountID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"message": "invalid token subject"}})
			return
		}

		if len(roleSet) > 0 {
			if _, ok := roleSet[role]; !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "error": gin.H{"message": "insufficient role"}})
				return
			}
		}

		c.Set(ContextAccountID, accountID)
		c.Set(ContextRole, role)
		c.Next()
	}
}
