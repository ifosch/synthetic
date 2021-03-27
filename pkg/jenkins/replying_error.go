package jenkins

// ReplyingError ...
type ReplyingError struct {
	Msg string
}

// Error ...
func (e ReplyingError) Error() string {
	return e.Msg
}
