package provider

import (
	"github.com/mattermost/mattermost-server/model"
	"testing"
		"reflect"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
)

func TestMattermost_NormalizeMessageEvent(t *testing.T) {
	post := &model.Post{
		Id: "post1",
		UserId: "user1",
		ChannelId: "channel1",
		Message: "Hello World",
	}

	botUser := &model.User{}
	botUser.Username = "Bot1"

	api := getNormalizeAPIMock()

	m := Mattermost {
		api,
		botUser,
	}
	expectedEvent := Event{
		Type: "message",
		PostID: "post1",
		UserID: "user1",
		UserName: "username1",
		UserRole: "system_user system_admin team_user team_admin",
		ChannelID: "channel1",
		ChannelName: "channelname1",
		ChannelType: "O",
		TeamID: "team1",
		TeamName: "teamname1",
		Text: "Hello World",
	}
	givenEvent := m.NormalizeMessageEvent(post)

	if !reflect.DeepEqual(expectedEvent, givenEvent) {
		t.Errorf("Expected %+v, got %+v", expectedEvent, givenEvent)
	}
}

func TestMattermost_NormalizeChannelJoinEvent(t *testing.T) {
	botUser := &model.User{}
	botUser.Username = "Bot1"

	api := getNormalizeAPIMock()

	m := Mattermost {
		api,
		botUser,
	}

	channelMember := &model.ChannelMember{
		UserId: "user1",
		ChannelId: "channel1",
	}

	expectedEvent := Event{
		Type: "channel_join",
		UserID: "user1",
		UserName: "username1",
		UserRole: "system_user system_admin team_user team_admin",
		ChannelID: "channel1",
		ChannelName: "channelname1",
		ChannelType: "O",
		TeamID: "team1",
		TeamName: "teamname1",
	}
	givenEvent := m.NormalizeChannelJoinEvent(channelMember, nil)

	if !reflect.DeepEqual(expectedEvent, givenEvent) {
		t.Errorf("Expected %+v, got %+v", expectedEvent, givenEvent)
	}

	actor := &model.User {
		Id: "actor1",
		Username: "actorname1",
		Roles: "system_user team_user",
	}

	expectedEvent.ActorID = "actor1"
	expectedEvent.ActorName = "actorname1"
	expectedEvent.ActorRole = "system_user team_user"

	givenEvent = m.NormalizeChannelJoinEvent(channelMember, actor)

	if !reflect.DeepEqual(expectedEvent, givenEvent) {
		t.Errorf("Expected %+v, got %+v", expectedEvent, givenEvent)
	}
}

func TestMattermost_NormalizeChannelLeaveEvent(t *testing.T) {
	botUser := &model.User{}
	botUser.Username = "Bot1"

	api := getNormalizeAPIMock()

	m := Mattermost {
		api,
		botUser,
	}

	channelMember := &model.ChannelMember{
		UserId: "user1",
		ChannelId: "channel1",
	}

	expectedEvent := Event{
		Type: "channel_leave",
		UserID: "user1",
		UserName: "username1",
		UserRole: "system_user system_admin team_user team_admin",
		ChannelID: "channel1",
		ChannelName: "channelname1",
		ChannelType: "O",
		TeamID: "team1",
		TeamName: "teamname1",
	}
	givenEvent := m.NormalizeChannelLeaveEvent(channelMember, nil)

	if !reflect.DeepEqual(expectedEvent, givenEvent) {
		t.Errorf("Expected %+v, got %+v", expectedEvent, givenEvent)
	}

	actor := &model.User {
		Id: "actor1",
		Username: "actorname1",
		Roles: "system_user team_user",
	}

	expectedEvent.ActorID = "actor1"
	expectedEvent.ActorName = "actorname1"
	expectedEvent.ActorRole = "system_user team_user"

	givenEvent = m.NormalizeChannelLeaveEvent(channelMember, actor)

	if !reflect.DeepEqual(expectedEvent, givenEvent) {
		t.Errorf("Expected %+v, got %+v", expectedEvent, givenEvent)
	}
}

func TestMattermost_NormalizeTeamJoinEvent(t *testing.T) {
	botUser := &model.User{}
	botUser.Username = "Bot1"

	api := getNormalizeAPIMock()

	m := Mattermost {
		api,
		botUser,
	}

	teamMember := &model.TeamMember{
		UserId: "user1",
		TeamId: "team1",
	}

	expectedEvent := Event{
		Type: "team_join",
		UserID: "user1",
		UserName: "username1",
		UserRole: "system_user system_admin team_user team_admin",
		TeamID: "team1",
		TeamName: "teamname1",
	}
	givenEvent := m.NormalizeTeamJoinEvent(teamMember, nil)

	if !reflect.DeepEqual(expectedEvent, givenEvent) {
		t.Errorf("Expected %+v, got %+v", expectedEvent, givenEvent)
	}

	actor := &model.User {
		Id: "actor1",
		Username: "actorname1",
		Roles: "system_user team_user",
	}

	expectedEvent.ActorID = "actor1"
	expectedEvent.ActorName = "actorname1"
	expectedEvent.ActorRole = "system_user team_user"

	givenEvent = m.NormalizeTeamJoinEvent(teamMember, actor)

	if !reflect.DeepEqual(expectedEvent, givenEvent) {
		t.Errorf("Expected %+v, got %+v", expectedEvent, givenEvent)
	}
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

func getNormalizeAPIMock() plugin.API{
	api := &plugintest.API{}
	api.On("GetChannel", "channel1").Return(&model.Channel{Id: "channel1", Name: "channelname1", TeamId: "team1", Type: "O"}, nil)
	api.On("GetUser", "user1").Return(&model.User{Id: "user1", Username: "username1", Roles:"system_user system_admin"}, nil)
	api.On("GetUser", "actor1").Return(&model.User{Id: "actor1", Username: "actorname1", Roles:"system_user"}, nil)
	api.On("GetTeamMember", "team1", "user1").Return(&model.TeamMember{UserId: "user1", TeamId: "team1", Roles:"team_user team_admin"}, nil)
	api.On("GetTeamMember", "team1", "actor1").Return(&model.TeamMember{UserId: "actor1", TeamId: "team1", Roles:"team_user"}, nil)
	api.On("GetTeam", "team1").Return(&model.Team{Id: "team1", Name: "teamname1"}, nil)
	return api
}