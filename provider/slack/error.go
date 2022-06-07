package slack

import "github.com/deepsourcelabs/hermes/domain"

type slackError struct {
	message    string
	internal   string
	statusCode int
	systemCode string
	isFatal    bool
}

func NewErr(statusCode int, systemCode, message, internal string, isFatal bool) domain.IError {
	return &slackError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}

func (e *slackError) Message() string {
	return e.message
}

func (e *slackError) IsFatal() bool {
	return e.isFatal
}

func (e *slackError) StatusCode() int {
	return e.statusCode
}

func (e *slackError) Error() string {
	return e.internal
}

func (e *slackError) SystemCode() string {
	return e.systemCode
}

var (
	errFailedOptsValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-SLK-40001", "unsupported notifier config", internal, true)
	}
	errFailedBodyValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-SLK-40001", "template incompatible for provider", internal, true)
	}
	errFailedSendTemporary = func(internal string) domain.IError {
		return NewErr(500, "ERR-SLK-50010", "failed to create issue", internal, false)
	}
	errFailedSendPermanent = func(internal string) domain.IError {
		return NewErr(500, "ERR-SLK-50020", "failed to create issue", internal, true)
	}
	errFailedOptsFetch = func(internal string) domain.IError {
		return NewErr(500, "ERR-SLK-50030", "failed to fetch opts", internal, true)
	}
)
