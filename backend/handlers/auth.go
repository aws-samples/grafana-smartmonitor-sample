package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"genai.ai/automonitor/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("e6b8e7a9c8e3b2f0d9e2f1d0c6b5a4b3a2")

func shouldSkipRewrite(path string) bool {
    skipPaths := []string{
        "/login",
        "/logout",
		"/static/",
        "/_next/",
        "/metrics",
		"/favicon.ico",
		"/scheduler",
		"/settings",
		"/",

    }

    for _, skipPath := range skipPaths {
        if path == skipPath || strings.HasPrefix(path, skipPath) {
            return true
        }
    }

    return false
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if shouldSkipRewrite(path) {
			c.Next()
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		
		// Check if the token string starts with "Bearer "
		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Extract the token string without the "Bearer " prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignInHandler(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ok := validateCredentials(creds.Username, creds.Password)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth failed"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": creds.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func validateCredentials(username, password string) bool {
	config := service.GetGlobalConfig()
	return username == "admin" && password == config.AdminPassword
}
