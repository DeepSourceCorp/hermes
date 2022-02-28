package domain

// IError is the base error interface used across the project.  This can be accepted as
// a standard `error` since it implements Error() method.
type IError interface {
	// Message returns the user world error message.  This must not leak
	// internal system details.
	Message() string

	// IsFatal defines if the error is transient and may be retried.
	IsFatal() bool

	// User world error code.  Use HTTPStatusCode for ease of use.
	StatusCode() int

	// Returns the actual error message.  This should be populated from the
	// source error.
	Error() string

	// System level error code.  These internal error codes can be
	// sent to the user world for pin pointing the error source.
	SystemCode() int
}
