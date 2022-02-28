package service

import (
	"github.com/deepsourcelabs/hermes/domain"
)

type serviceError struct {
	message    string
	internal   string
	statusCode int
	systemCode int
	isFatal    bool
}

func NewErr(statusCode int, systemCode int, message string, internal string, isFatal bool) domain.IError {
	return &serviceError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}
func (e *serviceError) Message() string {
	return e.message
}
func (e *serviceError) IsFatal() bool {
	return e.isFatal
}
func (e *serviceError) StatusCode() int {
	return e.statusCode
}
func (e *serviceError) Error() string {
	return e.message
}
func (e *serviceError) SystemCode() int {
	return e.systemCode
}

var (
	errMandatoryParamsMissing = func(internal string) domain.IError {
		return NewErr(400, 40001, "mandatory params missing", internal, true)
	}
	errUnprocessable = func(internal string) domain.IError {
		return NewErr(422, 40002, "unable to process the request", internal, true)
	}
	errRecipientMalformed = func(internal string) domain.IError {
		return NewErr(400, 40002, "some recipients are malformed", internal, true)
	}
	errMinOneRecipient = func(internal string) domain.IError {
		return NewErr(400, 40002, "at least one recipient must be defined", internal, true)
	}
	errStateless = func(internal string) domain.IError {
		return NewErr(422, 40002, "template lookup is not available", internal, true)
	}
)
