// Package exception — custom error types untuk error handling terpusat.
package exception

// NotFoundError — error ketika resource tidak ditemukan.
type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

// NewNotFoundError — konstruktor NotFoundError.
func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{Message: message}
}

// ValidationError — error ketika input tidak valid.
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// NewValidationError — konstruktor ValidationError.
func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}

// ForbiddenError — error ketika user tidak punya akses.
type ForbiddenError struct {
	Message string
}

func (e ForbiddenError) Error() string {
	return e.Message
}
