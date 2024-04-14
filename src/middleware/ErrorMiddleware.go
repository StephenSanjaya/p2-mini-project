package middleware

import (
	"p2-mini-project/src/httputil"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware(c *gin.Context) {
	c.Next()

	err := c.Errors.Last()
	if err != nil {
		switch e := err.Err.(type) {
		case *httputil.HTTPError:
			c.JSON(e.Code, gin.H{
				"message": e.Message,
				"detail":  e.Detail,
			})
		default:
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
		c.Abort()
	}
}
