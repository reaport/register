package errors

import (
	"errors"
)

var (
	ErrValid          = errors.New("validation errors")
	ErrBaggageSize    = errors.New("exceeding the allowed baggage size")
	ErrTicketNotFound = errors.New("the passenger was not found for registration")
	ErrInternalServer = errors.New("internal server errors")
)

func GetCode(message string) int {
	switch message {
	case ErrValid.Error():
		return 400
	case ErrBaggageSize.Error():
		return 400
	case ErrTicketNotFound.Error():
		return 404
	case ErrInternalServer.Error():
		return 500
	default:
		return 500
	}
}
