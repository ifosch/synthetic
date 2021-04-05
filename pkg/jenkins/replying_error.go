package jenkins

// ReplyingError is an error type to identify errors to be replied to
// user.
type ReplyingError struct {
	Msg string
}

// Error satisfy the `error` interface.
func (e ReplyingError) Error() string {
	return e.Msg
}
