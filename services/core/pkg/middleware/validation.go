package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/papidb/aqua/pkg/api"
)

// ErrNotJSON is the error when the request body is not JSON
var ErrNotJSON = errors.New("body is not JSON")

// ValidationMiddleware is a middleware to validate incoming requests
func ValidationMiddleware(v validation.Validatable) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the request body for logging purposes (before binding)
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}

		// Reassign the body so it can be used by the next handler (e.g., binding)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Check Content-Type header for JSON
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": ErrNotJSON.Error()})
			c.Abort()
			return
		}

		// Decode the body into the struct
		err = json.NewDecoder(c.Request.Body).Decode(v)
		if err != nil && err != io.EOF {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
			c.Abort()
			return
		}

		// Perform validation using ozzo-validation
		err = validation.Validate(v)

		// If validation passes, proceed to the next middleware or handler
		if err == nil {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			c.Next()
			return
		}
		// If validation fails, return the validation errors
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
				Data:    map[string]interface{}{"validation_error": e},
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
		}

	}
}
