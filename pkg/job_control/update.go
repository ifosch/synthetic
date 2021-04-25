package jobcontrol

// Update is a message update.
type Update struct {
	Msg      string
	Reaction string
	Done     bool
}
