package github

import "github.com/deepsourcelabs/hermes/domain"

type githubError struct {
	message    string
	internal   string
	statusCode int
	systemCode string
	isFatal    bool
}

func NewErr(statusCode int, systemCode, message, internal string, isFatal bool) domain.IError {
	return &githubError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}

func (e *githubError) Message() string {
	return e.message
}
func (e *githubError) IsFatal() bool {
	return e.isFatal
}
func (e *githubError) StatusCode() int {
	return e.statusCode
}
func (e *githubError) Error() string {
	return e.internal
}
func (e *githubError) SystemCode() string {
	return e.systemCode
}

var (
	errFailedOptsValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-GTH-40001", "unsupported notifier config", internal, true)
	}
	errFailedBodyValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-GTH-40001", "template incompatible for provider", internal, true)
	}
	errFailedSendTemporary = func(internal string) domain.IError {
		return NewErr(500, "ERR-GTH-50010", "failed to create issue", internal, false)
	}
	errFailedSendPermanent = func(internal string) domain.IError {
		return NewErr(500, "ERR-GTH-50020", "failed to create issue", internal, true)
	}
)
