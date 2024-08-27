package infrastructure

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return

		}

		authParts := strings.Split(authHeader, " ")

		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return

		}

		token, err := jwt.Parse(authParts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return "Secret", nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid JWT", "err": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Accessing nested claims correctly
			if customClaims, ok := claims["claims"].(map[string]interface{}); ok {
				userID, userIDOk := customClaims["user_id"].(string)
				role, roleOk := customClaims["role"].(string)
				email, emailOk := customClaims["email"].(string)

				if userIDOk && roleOk && emailOk {
					c.Set("user_id", userID)
					c.Set("role", role)
					c.Set("email", email)

				} else {
					c.JSON(401, gin.H{"error": "Invalid JWT claims"})
					c.Abort()
					return
				}
			} else {
				c.JSON(401, gin.H{"error": "Invalid JWT claims structure"})
				c.Abort()
				return
			}
		} else {
			c.JSON(401, gin.H{"error": "Invalid JWT claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}
