package httputil

import "github.com/gin-gonic/gin"

type HTTPSuccess struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewSuccess(c *gin.Context, status int, message string, data ...any) {
	su := HTTPSuccess{
		Message: message,
		Data:    data,
	}
	c.JSON(status, su)
}
