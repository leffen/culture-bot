package slack

import (
	"os"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

// New slack obj
func New() *slack.Client {
	return slack.New(os.Getenv("SLACK_TOKEN"))
}

// GetAllUsers returns all users with email
func GetAllUsers(c *slack.Client) map[string]slack.User {
	users, err := c.GetUsers()
	if err != nil {
		log.WithError(err).Error("unable to get slack users")
		return nil
	}

	activeUsersWithEmail := make(map[string]slack.User)
	for _, u := range users {
		if !u.IsBot && !u.Deleted && len(u.Profile.Email) > 0 {
			if len(u.Profile.DisplayName) > 0 {
				activeUsersWithEmail[u.Profile.DisplayName] = u
			} else {
				activeUsersWithEmail[u.Name] = u
			}
		}
	}

	return activeUsersWithEmail
}
