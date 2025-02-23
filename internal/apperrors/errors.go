package apperrors

import "errors"

var (
	ErrFlightNotFound = errors.New("❌ flight not found")
	ErrSeatNotFound   = errors.New("❌ not found place ")
	ErrSeatTaken      = errors.New("❌ seat already taken")
	ErrTicketNotFound = errors.New("❌ ticket not found")
	ErrInternalServer = errors.New("❌internal server error")
)
