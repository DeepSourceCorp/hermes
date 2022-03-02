package jira

import "github.com/deepsourcelabs/hermes/domain"

type jiraError struct {
	message    string
	internal   string
	systemCode string
	statusCode int
	isFatal    bool
}

func NewErr(statusCode int, systemCode string, message string, internal string, isFatal bool) domain.IError {
	return &jiraError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}
func (e *jiraError) Message() string {
	return e.message
}
func (e *jiraError) IsFatal() bool {
	return e.isFatal
}
func (e *jiraError) StatusCode() int {
	return e.statusCode
}
func (e *jiraError) Error() string {
	return e.internal
}
func (e *jiraError) SystemCode() string {
	return e.systemCode
}

var (
	errFailedOptsValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-JIR-40001", "unsupported notifier config", internal, true)
	}
	errFailedBodyValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-JIR-40001", "template incompatible for provider", internal, true)
	}
	errFailedSendTemporary = func(internal string) domain.IError {
		return NewErr(500, "ERR-JIR-50010", "failed to create issue", internal, false)
	}
	errFailedSendPermanent = func(internal string) domain.IError {
		return NewErr(500, "ERR-JIR-50020", "failed to create issue", internal, true)
	}
)
