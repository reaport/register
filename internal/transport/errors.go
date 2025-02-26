package transport

const (
	ErrValidation = "Validation errors"
	ErrInternal   = "Internal server errors"
)

type ErrorResponse struct {
	Message string `json:"errors"`
}
