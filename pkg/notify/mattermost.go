package notify

import (
	"fmt"
	"net/http"

	mattermost "github.com/mattermost/mattermost-server/v5/model"
)

type Mattermost struct {
	Client    *mattermost.Client4
	ChannelID string
}

// NewMattermost returns an instance of an authenticated
// *mattermost.Client4 and the mattermost direct message channel ID
func NewMattermost(url, token, username string) (Notifier, error) {
	client := mattermost.NewAPIv4Client(url)
	client.AuthToken = token
	client.AuthType = mattermost.HEADER_AUTH
	client.HttpHeader = map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	botUser, res := client.GetMe("")
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to authenticate bot user: %s, using token", botUser.Nickname)
	}
	sendUser, res := client.GetUserByUsername(username, "")
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to find user id of user: %s", username)
	}
	directChannel, res := client.CreateDirectChannel(botUser.Id, sendUser.Id)
	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unable to crete direct message between botUser: %s and %s", botUser.Nickname, username)
	}
	return &Mattermost{
		Client:    client,
		ChannelID: directChannel.Id,
	}, nil
}

// SendMessage sends the message to a mattermost user as a bot
func (m *Mattermost) SendMessage(body string) error {
	if _, res := m.Client.CreatePost(&mattermost.Post{
		ChannelId: m.ChannelID,
		Message:   body,
	}); res.StatusCode != http.StatusCreated {
		return fmt.Errorf("error sending message to channel: %s", m.ChannelID)
	}
	return nil
}
