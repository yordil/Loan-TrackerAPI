package infrastructure

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(secret string) gin.HandlerFunc {

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

			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid JWT", "err": err.Error()})
			c.Abort()
			return
		}


		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, userIDOk := claims["id"].(string)
			if !userIDOk {
				c.JSON(401, gin.H{"error": "Invalid JWT claims"})
				c.Abort()
				return
			}

			// Store the user ID in the context
			c.Set("user_id", userID)
		} else {
			c.JSON(401, gin.H{"error": "Invalid JWT claims"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}
