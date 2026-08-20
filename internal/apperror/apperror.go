package apperror

import "net/http"

// AppError is a typed error that carries an HTTP status code alongside a
// user-safe message. Services and repositories return these instead of
// writing gin.JSON() themselves — that keeps them free of any HTTP import,
// which is what makes them unit-testable and reusable outside Gin.
type AppError struct {
	Code    int    // HTTP status code this error maps to
	Message string // safe to show the client
	Err     error  // wrapped internal error, for logging only — never sent to the client
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap lets errors.Is / errors.As see through to the wrapped error.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates an AppError with no wrapped internal error.
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap creates an AppError around an internal error (e.g. a GORM error),
// keeping the internal detail out of the client-facing Message.
func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// Common constructors — use these in services/repositories instead of
// reaching for net/http status codes directly, so the mapping stays
// consistent everywhere.
func BadRequest(message string) *AppError   { return New(http.StatusBadRequest, message) }
func Unauthorized(message string) *AppError { return New(http.StatusUnauthorized, message) }
func Forbidden(message string) *AppError    { return New(http.StatusForbidden, message) }
func NotFound(message string) *AppError     { return New(http.StatusNotFound, message) }
func Conflict(message string) *AppError     { return New(http.StatusConflict, message) }

func Internal(message string, err error) *AppError {
	return Wrap(http.StatusInternalServerError, message, err)
}
