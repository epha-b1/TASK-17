package server

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func abortAPIError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, apiError{Code: code, Message: message})
}

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func isValidUUID(s string) bool {
	return uuidRegex.MatchString(s)
}
