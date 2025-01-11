package api

import "fmt"

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
