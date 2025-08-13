package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Authentication - JWT authentication
func Authentication() gin.HandlerFunc {
	// Load .env and secret key
	godotenv.Load()
	jwtSecret := []byte(strings.TrimSpace(os.Getenv("JWT_SECRET_KEY")))

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// // Skip auth for Swagger/OpenAPI docs
		// if strings.HasPrefix(c.Request.URL.Path, "/docs") ||
		// 	c.Request.URL.Path == "/openapi.json" ||
		// 	c.Request.URL.Path == "/swagger.json" {
		// 	c.Next()
		// 	return
		// }

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Authorization token is missing"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid authorization header format."})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid token."})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"detail": "Token has expired."})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"detail": fmt.Sprintf("Invalid token: %v", err)})
			}
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Expiration check
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				c.JSON(http.StatusUnauthorized, gin.H{"detail": "Token has expired."})
				c.Abort()
				return
			}
		}

		// Email check
		email := fmt.Sprintf("%v", claims["email"])
		if email == "" || email == "<nil>" {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid token: missing email."})
			c.Abort()
			return
		}

		// Issuer check
		issuer := fmt.Sprintf("%v", claims["iss"])
		if issuer == "" || issuer == "<nil>" {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid token: missing issuer."})
			c.Abort()
			return
		}

		// Attach email & issuer to context
		c.Set("email", email)
		c.Set("issuer", issuer)

		c.Next()
	}
}
