package web

const (
	ValidationErr    = "Validation"
	SerializationErr = "Serialization"
	NotFoundErr      = "NotFound"
)

// validation error
type ValidationError struct {
	Type    string      `json:"-"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func (e *ValidationError) Error() string {
	return ValidationErr + " " + e.Message
}

type SerializationError struct {
	Type    string `json:"-"`
	Message string `json:"message"`
}

func (e *SerializationError) Error() string {
	return SerializationErr + " " + e.Message
}

type NotFoundError struct {
	Type    string `json:"-"`
	Message string `json:"message"`
}

func (e *NotFoundError) Error() string {
	return NotFoundErr + " " + e.Message
}
