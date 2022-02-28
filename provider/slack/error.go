package slack

import "github.com/deepsourcelabs/hermes/domain"

type slackError struct {
	message    string
	internal   string
	statusCode int
	systemCode int
	isFatal    bool
}

func NewErr(statusCode int, systemCode int, message string, internal string, isFatal bool) domain.IError {
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
func (e *slackError) SystemCode() int {
	return e.systemCode
}

var (
	errSlackErr = func(internal string) domain.IError {
		return NewErr(500, 50001, "failed to send message to Slack", internal, true)
	}
	errSlackOptsParseFail = func(internal string) domain.IError {
		return NewErr(400, 40010, "failed to parse Slack options", internal, true)
	}
)
