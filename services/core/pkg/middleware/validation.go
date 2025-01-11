package middlewares

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/papidb/aqua/pkg/api"

	"github.com/gin-gonic/gin"
)

var ErrNotJSON = errors.New("body is not JSON")

// ValidationMiddleware is a middleware to validate incoming requests
func ValidationMiddleware(v validation.Validatable) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the request JSON to the provided struct

		r := c.Request

		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": ErrNotJSON.Error()})
			c.Abort()
			return
		}

		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil && err != io.EOF {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
			c.Abort()
			return
		}

		err = validation.Validate(v)
		if err == nil {
			return
		}

		var e validation.Errors
		switch {
		case err == ErrNotJSON:
			c.JSON(http.StatusUnsupportedMediaType, &api.AppErr{
				Code:    http.StatusUnsupportedMediaType,
				Message: http.StatusText(http.StatusUnsupportedMediaType),
				Err:     err,
			})
			c.Abort()
			return
		case errors.As(err, &e):
			c.JSON(http.StatusBadRequest, &api.AppErr{
				Code:    http.StatusBadRequest,
				Message: "We could not validate your request.",
				Data:    map[string]any{"validation_error": e},
			})
			c.Abort()
			return
		default:
			c.JSON(http.StatusBadRequest, &api.AppErr{
				Code:    http.StatusBadRequest,
				Message: "We cannot parse your request body.",
				Err:     err,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
