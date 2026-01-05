package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Gauravsulegai/careersync/internal/database"
	"github.com/Gauravsulegai/careersync/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// 1. Get the Authorization header
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// 2. The header usually looks like "Bearer eyJhbGci..."
	// We need to split it to get just the token part
	tokenString := strings.Split(authHeader, " ")
	if len(tokenString) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// 3. Parse and Validate the token
	token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("my_secret_key"), nil // MUST match the key used in auth.go
	})

	if err != nil {
		// This will tell us if the key is wrong or the format is bad
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is claimed to be invalid"})
		return
	}

	// 4. Extract user data from token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}

		// Find the user in the DB
		var user models.User
		database.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// 5. Attach user to the request (So we can use it later)
		c.Set("user", user)
		
		// Move to the next handler
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
	}
}