package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"p2-mini-project/src/httputil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(roles string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("authorization")

		if tokenString == "" {
			c.Error(httputil.NewError(http.StatusUnauthorized, "unauthorized", errors.New("token not found")))
			c.Abort()
			return
		}

		secret_token := []byte(os.Getenv("JWT"))
		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.Error(httputil.NewError(http.StatusUnauthorized, "unauthorized", errors.New("invalid algorithm use")))
				c.Abort()
				return nil, fmt.Errorf("invalid algorithm use")
			}
			return secret_token, nil
		})

		if parsedToken == nil || !parsedToken.Valid {
			c.Error(httputil.NewError(http.StatusUnauthorized, "unauthorized", err))
			c.Abort()
			return
		}

		if float64(time.Now().Unix()) > parsedToken.Claims.(jwt.MapClaims)["exp"].(float64) {
			c.Error(httputil.NewError(http.StatusUnauthorized, "unauthorized", errors.New("token expired")))
			c.Abort()
			return
		}

		user_role := parsedToken.Claims.(jwt.MapClaims)["role"]
		if user_role == "user" && roles == "admin" {
			c.Error(httputil.NewError(http.StatusUnauthorized, "unauthorized", errors.New("need admin role to access this api")))
			c.Abort()
			return
		}

		c.Set("user_id", parsedToken.Claims.(jwt.MapClaims)["user_id"])

		c.Next()
	}
}
