package discord

import "github.com/deepsourcelabs/hermes/domain"

type discordError struct {
	message    string
	internal   string
	statusCode int
	systemCode string
	isFatal    bool
}

func NewErr(statusCode int, systemCode, message, internal string, isFatal bool) domain.IError {
	return &discordError{
		message:    message,
		statusCode: statusCode,
		systemCode: systemCode,
		internal:   internal,
		isFatal:    isFatal,
	}
}

func (e *discordError) Message() string {
	return e.message
}
func (e *discordError) IsFatal() bool {
	return e.isFatal
}
func (e *discordError) StatusCode() int {
	return e.statusCode
}
func (e *discordError) Error() string {
	return e.internal
}
func (e *discordError) SystemCode() string {
	return e.systemCode
}

var (
	errFailedOptsValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-DIS-40001", "unsupported notifier config", internal, true)
	}
	errFailedBodyValidation = func(internal string) domain.IError {
		return NewErr(422, "ERR-DIS-40001", "template incompatible for provider", internal, true)
	}
	errFailedSendTemporary = func(internal string) domain.IError {
		return NewErr(500, "ERR-DIS-50010", "failed to send message", internal, false)
	}
	errFailedSendPermanent = func(internal string) domain.IError {
		return NewErr(500, "ERR-DIS-50020", "failed to send message", internal, true)
	}
)
