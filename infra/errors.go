package infra

type AppError struct {
	Code    int    `json:",omitempty"`
	Message string `json:"message"`
}

func (e AppError) AsMessage() *AppError {
	return &AppError{Message: e.Message}
}

func NewError(message string, code int) *AppError {
	return &AppError{Message: message, Code: code}
}
