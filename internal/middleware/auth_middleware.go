package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hunafazaky/event-booking-app/internal/response"
)

func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, "invalid or missing token")
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || token == nil || !token.Valid {
			response.Fail(c, http.StatusUnauthorized, "invalid or missing token")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(float64); ok {
				c.Set("user_id", uint(sub))
				c.Next()
				return
			}
		}

		response.Fail(c, http.StatusUnauthorized, "invalid or missing token")
		c.Abort()
	}
}
