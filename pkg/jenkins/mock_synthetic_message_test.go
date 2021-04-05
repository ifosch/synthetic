package jenkins

import "github.com/ifosch/synthetic/pkg/synthetic"

type MockSyntheticUser struct {
	name string
}

func (msu MockSyntheticUser) Name() string {
	return msu.name
}

type MockSyntheticConversation struct {
	name string
}

func (msc MockSyntheticConversation) Name() string {
	return msc.name
}

type MockSyntheticMessage struct {
	thread       bool
	mention      bool
	text         string
	user         MockSyntheticUser
	conversation MockSyntheticConversation
	replies      []string
}

func (msm *MockSyntheticMessage) Reply(msg string, inThread bool) {
	msm.replies = append(msm.replies, msg)
}

func (msm *MockSyntheticMessage) React(reaction string) {
}

func (msm *MockSyntheticMessage) Unreact(reaction string) {
}

func (msm *MockSyntheticMessage) ClearMention() string {
	return msm.text
}

func (msm *MockSyntheticMessage) Thread() bool {
	return msm.thread
}

func (msm *MockSyntheticMessage) Mention() bool {
	return msm.mention
}

func (msm *MockSyntheticMessage) Text() string {
	return msm.text
}

func (msm *MockSyntheticMessage) User() synthetic.User {
	return msm.user
}

func (msm *MockSyntheticMessage) Conversation() synthetic.Conversation {
	return msm.conversation
}
