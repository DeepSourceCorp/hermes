package mailgun

import "github.com/deepsourcelabs/hermes/domain"

type mailgunError struct {
	message    string
	internal   string
	statusCode int
	systemCode string
	isFatal    bool
}

func NewErr(statusCode int, systemCode, message, internal string, isFatal bool) domain.IError {
	return &mailgunError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}

func (e *mailgunError) Message() string {
	return e.message
}
func (e *mailgunError) IsFatal() bool {
	return e.isFatal
}
func (e *mailgunError) StatusCode() int {
	return e.statusCode
}
func (e *mailgunError) Error() string {
	return e.internal
}
func (e *mailgunError) SystemCode() string {
	return e.systemCode
}

var (
	errFailedOptsValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-MLG-40001", "unsupported notifier config", internal, true)
	}
	errFailedBodyValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-MLG-40001", "template incompatible for provider", internal, true)
	}
	errFailedSendTemporary = func(internal string) domain.IError {
		return NewErr(500, "ERR-MLG-50010", "failed to send mail", internal, false)
	}
	errFailedSendPermanent = func(internal string) domain.IError {
		return NewErr(500, "ERR-MLG-50020", "failed to send mail", internal, true)
	}
)
