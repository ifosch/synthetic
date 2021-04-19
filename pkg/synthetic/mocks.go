package synthetic

// MockUser is a mock of a User.
type MockUser struct {
	name string
}

// Name is a mock for User.Name() method.
func (msu MockUser) Name() string {
	return msu.name
}

// MockConversation is a mock for a Conversation
type MockConversation struct {
	name string
}

// Name is a mock for Conversation.Name() method.
func (msc MockConversation) Name() string {
	return msc.name
}

// MockMessage is a mock for a Message.
type MockMessage struct {
	thread       bool
	mention      bool
	text         string
	user         MockUser
	conversation MockConversation
	replies      []string
}

// NewMockMessage is the MockMessage constructor.
func NewMockMessage(input string) *MockMessage {
	return &MockMessage{
		text:    input,
		replies: []string{},
	}
}

// Replies returns the replies received by the MockMessage.
func (msm *MockMessage) Replies() []string {
	return msm.replies
}

// Reply is a mock for Message.Reply() method.
func (msm *MockMessage) Reply(msg string, inThread bool) {
	msm.replies = append(msm.replies, msg)
}

// React is a mock for Message.React() method.
func (msm *MockMessage) React(reaction string) {
}

// Unreact is a mock for Message.Unreact() method.
func (msm *MockMessage) Unreact(reaction string) {
}

// ClearMention is a mock for Message.ClearMention() method.
func (msm *MockMessage) ClearMention() string {
	return msm.text
}

// Thread is a mock for Message.Thread() method.
func (msm *MockMessage) Thread() bool {
	return msm.thread
}

// Mention is a mock for Message.Mention() method.
func (msm *MockMessage) Mention() bool {
	return msm.mention
}

// Text is a mock for Message.Text() method.
func (msm *MockMessage) Text() string {
	return msm.text
}

// User is a mock for Message.User() method.
func (msm *MockMessage) User() User {
	return msm.user
}

// Conversation is a mock for Message.Conversation() method.
func (msm *MockMessage) Conversation() Conversation {
	return msm.conversation
}
