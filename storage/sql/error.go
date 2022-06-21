package storage

import "github.com/deepsourcelabs/hermes/domain"

type sqlError struct {
	message    string
	internal   string
	statusCode int
	systemCode string
	isFatal    bool
}

func NewErr(statusCode int, systemCode, message, internal string, isFatal bool) domain.IError {
	return &sqlError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}

func (e *sqlError) Message() string {
	return e.message
}

func (e *sqlError) IsFatal() bool {
	return e.isFatal
}

func (e *sqlError) StatusCode() int {
	return e.statusCode
}

func (e *sqlError) Error() string {
	return e.internal
}

func (e *sqlError) SystemCode() string {
	return e.systemCode
}

var errDBErr = func(internal string) domain.IError {
	return NewErr(500, "HE-STO-50010", "something went wrong", internal, true)
}
