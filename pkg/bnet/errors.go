package bnet

import (
	"errors"
	"fmt"
)

var ErrTokenIsInvalid error = errors.New("the provided token is invalid")

type ErrMissingRequiredScope struct {
	Scope string
}

func (m ErrMissingRequiredScope) Error() string {
	return fmt.Sprintf("missing the required scope '%s'", m.Scope)
}
