package transport

const (
	ErrValidation = "Validation error"
	ErrInternal   = "Internal server error"
)

type ErrorResponse struct {
	Message string `json:"error"`
}
