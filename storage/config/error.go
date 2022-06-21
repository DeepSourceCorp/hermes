package config

import "github.com/deepsourcelabs/hermes/domain"

type cfgError struct {
	message    string
	internal   string
	statusCode int
	systemCode string
	isFatal    bool
}

func NewErr(statusCode int, systemCode, message, internal string, isFatal bool) domain.IError {
	return &cfgError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}

func (e *cfgError) Message() string {
	return e.message
}

func (e *cfgError) IsFatal() bool {
	return e.isFatal
}

func (e *cfgError) StatusCode() int {
	return e.statusCode
}

func (e *cfgError) Error() string {
	return e.internal
}

func (e *cfgError) SystemCode() string {
	return e.systemCode
}

var errDBErr = func(internal string) domain.IError {
	return NewErr(500, "HE-STO-50010", "something went wrong", internal, true)
}
