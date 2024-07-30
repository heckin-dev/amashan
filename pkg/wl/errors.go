package wl

import "fmt"

type ErrNoPointsLeft struct {
	StatusCode       int
	RemainingSeconds int
	Err              error
}

func (e *ErrNoPointsLeft) Error() string {
	return fmt.Sprintf("the service is currently unavailable, try again in '%d' seconds", e.RemainingSeconds)
}

func (e *ErrNoPointsLeft) Unwrap() error {
	return e.Err
}
