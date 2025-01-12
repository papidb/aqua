package api

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// error not found
var ErrCustomerNotFound = errors.New("customer not found")
var ErrResourceNotFound = errors.New("resource not found")

type AppErr struct {
	Code    int         `json:"-"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Err     error       `json:"-"`
}

func (e AppErr) Error() string {
	if e.Err == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e AppErr) Unwrap() error { return e.Err }

func HandleMappedErrors(c *gin.Context, err error, errorMapping map[error]int) bool {
	for e, status := range errorMapping {
		if errors.Is(err, e) {
			Error(c.Request, c.Writer, AppErr{
				Code:    status,
				Message: err.Error(),
				Err:     err,
			})
			return true
		}
	}
	return false
}
