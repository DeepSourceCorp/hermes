package smtp

import "github.com/deepsourcelabs/hermes/domain"

type smtpError struct {
	message    string
	internal   string
	statusCode int
	systemCode string
	isFatal    bool
}

func NewErr(statusCode int, systemCode, message, internal string, isFatal bool) domain.IError {
	return &smtpError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}

func (e *smtpError) Message() string {
	return e.message
}
func (e *smtpError) IsFatal() bool {
	return e.isFatal
}
func (e *smtpError) StatusCode() int {
	return e.statusCode
}
func (e *smtpError) Error() string {
	return e.internal
}
func (e *smtpError) SystemCode() string {
	return e.systemCode
}

var (
	errFailedOptsValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-SMT-40001", "unsupported notifier config", internal, true)
	}
	errFailedBodyValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-SMT-40001", "template incompatible for provider", internal, true)
	}
	errFailedSendTemporary = func(internal string) domain.IError {
		return NewErr(500, "ERR-SMT-50010", "failed to send mail", internal, false)
	}
)
