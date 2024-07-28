package bnet

import (
	"errors"
	"fmt"
)

var (
	ErrTokenIsInvalid = errors.New("the provided token is invalid")
)

type ErrMissingRequiredScope struct {
	Scope string
}

func (m ErrMissingRequiredScope) Error() string {
	return fmt.Sprintf("missing the required scope '%s'", m.Scope)
}

type ErrUnexpectedResponse struct {
	StatusCode int
}

func (e ErrUnexpectedResponse) Error() string {
	return fmt.Sprintf("unexpected response from server with status '%d'", e.StatusCode)
}
