package models

import (
	"fmt"
)

// A error that we can use to throw client-ready error messages
type RouteError struct {
	Message RouteErrorMessage `json:"error"`
	Code    int
}

type RouteErrorMessage struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (ae *RouteError) Error() string {
	return fmt.Sprintf("%s(%d): %s", ae.Message.Title, ae.Code, ae.Message.Description)
}

func (ae *RouteError) Is(target error) bool {
	_, isRE := target.(*RouteError)
	return isRE
}
