package storage

import "github.com/deepsourcelabs/hermes/domain"

type storageError struct {
	message    string
	internal   string
	statusCode int
	systemCode int
	isFatal    bool
}

func NewErr(statusCode int, systemCode int, message string, internal string, isFatal bool) domain.IError {
	return &storageError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}
func (e *storageError) Message() string {
	return e.message
}
func (e *storageError) IsFatal() bool {
	return e.isFatal
}
func (e *storageError) StatusCode() int {
	return e.statusCode
}
func (e *storageError) Error() string {
	return e.internal
}
func (e *storageError) SystemCode() int {
	return e.systemCode
}

var (
	errDBErr = func(internal string) domain.IError {
		return NewErr(500, 50001, "something went wrong", internal, true)
	}
)
