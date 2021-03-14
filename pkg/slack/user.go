package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

// User is ...
type User struct {
	slackUser *slack.User
	Name      string
}

// NewUserFromID ...
func NewUserFromID(id string, api *slack.Client) (user *User, err error) {
	userInfo, err := api.GetUserInfo(id)
	if err != nil {
		return nil, err
	}
	user = &User{userInfo, fmt.Sprintf("@%v", userInfo.Name)}
	return user, err
}
