package middlewares

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"golang-restapi/config"
	"golang-restapi/utils"
)

var jwtSecret []byte

func init() {
	godotenv.Load() // Load environment variables from .env
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

// Authentication verifies JWT token and checks expiration & email
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Check the Token string is in the header
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims (email & expiration)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Check expiration
		exp, ok := claims["exp"].(float64)
		if !ok || int64(exp) < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// Check email exists in token
		email, ok := claims["email"].(string)
		if !ok || email == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
			c.Abort()
			return
		}

		cfg := config.LoadConfig()
		issuer, ok := claims["iss"].(string)
		if !ok || issuer == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Issuer(iss) not found in token"})
			c.Abort()
			return
		}

		println("cfg.AuthIssuers:", cfg.AuthIssuers, issuer, err)

		// check the issuer is in the list
		if !utils.IsInList(cfg.AuthIssuers, issuer) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid issuer"})
			c.Abort()
			return
		}

		/*
		   // check the email is in the list, maybe does not need and issuer will be enough
		   if !utils.IsInList(cfg.AuthEmails, email) {
		       c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		       c.Abort()
		       return
		   }
		*/

		// Store email,issuer in Gin context for later use
		c.Set("email", email)
		c.Set("issuer", issuer)
		c.Next()
	}
}
