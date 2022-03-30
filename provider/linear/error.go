package linear

import "github.com/deepsourcelabs/hermes/domain"

type linearError struct {
	message    string
	internal   string
	statusCode int
	systemCode string
	isFatal    bool
}

func NewErr(statusCode int, systemCode, message, internal string, isFatal bool) domain.IError {
	return &linearError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}

func (e *linearError) Message() string {
	return e.message
}
func (e *linearError) IsFatal() bool {
	return e.isFatal
}
func (e *linearError) StatusCode() int {
	return e.statusCode
}
func (e *linearError) Error() string {
	return e.internal
}
func (e *linearError) SystemCode() string {
	return e.systemCode
}

var (
	errFailedOptsValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-LIN-40001", "unsupported notifier config", internal, true)
	}
	errFailedBodyValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-LIN-40001", "template incompatible for provider", internal, true)
	}
	errFailedSendTemporary = func(internal string) domain.IError {
		return NewErr(500, "ERR-LIN-50010", "failed to create issue", internal, false)
	}
	errFailedSendPermanent = func(internal string) domain.IError {
		return NewErr(500, "ERR-LIN-50020", "failed to create issue", internal, true)
	}
)
