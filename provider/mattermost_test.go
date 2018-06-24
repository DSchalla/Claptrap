package provider

import (
	"github.com/mattermost/mattermost-server/model"
	"testing"
)

func TestMattermost_handleMessageEvent(t *testing.T) {

}

func TestMattermost_ReplaceMessagePlaceholders(t *testing.T) {
	botUser := model.User{}
	botUser.Username = "Bot1"

	event := Event{}
	event.UserName = "User1"
	event.ActorName = "Actor1"
	event.ChannelName = "Channel1"

	message := "{User_Name}-{Actor_Name}-{Channel_Name}-{Bot_Name}"
	expected := "User1-Actor1-Channel1-Bot1"

	m := Mattermost{
		nil,
		&botUser,
	}
	given := m.ReplaceMessagePlaceholders(event, message)

	if expected != given {
		t.Errorf("expected %s, got %s", expected, given)
	}
}
