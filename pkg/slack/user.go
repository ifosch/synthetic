package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

// User is a weapper over slack-go's User object. It provides some
// utility methods over the User information.
type User struct {
	slackUser *slack.User
	name      string
}

// NewUserFromID returns a User object wrapping the user identified by
// `id`.
func NewUserFromID(id string, api IClient) (user *User, err error) {
	userInfo, err := api.GetUserInfo(id)
	if err != nil {
		return nil, err
	}
	user = &User{userInfo, fmt.Sprintf("@%v", userInfo.Name)}
	return user, err
}

// Name returns the name of the user.
func (u *User) Name() string {
	return u.name
}
