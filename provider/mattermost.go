package provider

import (
	"github.com/mattermost/mattermost-server/model"
	"log"
	"strings"
	"github.com/mattermost/mattermost-server/plugin"
)

type Mattermost struct {
	api     plugin.API
	botUser *model.User
}

func NewMattermost(api plugin.API, botUser *model.User) *Mattermost {
	m := Mattermost{}
	m.api = api
	m.botUser = botUser
	return &m
}

func (m *Mattermost) AutoJoinAllChannel() error {
	teams, err := m.api.GetTeams()

	if err != nil {
		return err
	}

	for _, team := range teams {
		_, err := m.api.GetPublicChannelsForTeam(team.Id, 0, 500)

		if err != nil {
			return err
		}

	}

	return nil
}

func (m *Mattermost) NormalizeMessageEvent(post *model.Post) Event {

	if post.Type == "system_add_to_channel" {
		return m.NormalizeUserInviteEvent(post)
	} else if post.Type == "system_join_channel" {
		return m.NormalizeUserJoinEvent(post)
	}

	unifiedEvent := Event{}
	unifiedEvent.Type = "message"
	unifiedEvent.PostID = post.Id
	unifiedEvent.UserID = post.UserId
	unifiedEvent.ChannelID = post.ChannelId
	unifiedEvent.Timestamp = post.CreateAt
	unifiedEvent.Text = post.Message
	unifiedEvent = m.addEventMetadata(unifiedEvent)
	return unifiedEvent
}

func (m *Mattermost) NormalizeUserInviteEvent(post *model.Post) Event {
	unifiedEvent := Event{}
	unifiedEvent.Type = "user_add"
	unifiedEvent.PostID = post.Id
	unifiedEvent.ChannelID = post.ChannelId
	unifiedEvent.UserName = post.Props["addedUsername"].(string)
	unifiedEvent.ActorName = post.Props["username"].(string)
	unifiedEvent = m.addEventMetadata(unifiedEvent)
	return unifiedEvent
}

func (m *Mattermost) NormalizeUserJoinEvent(post *model.Post) Event {
	unifiedEvent := Event{}
	unifiedEvent.Type = "user_add"
	unifiedEvent.PostID = post.Id
	unifiedEvent.ChannelID = post.ChannelId
	unifiedEvent.UserName = post.Props["username"].(string)
	unifiedEvent = m.addEventMetadata(unifiedEvent)
	return unifiedEvent
}

func (m *Mattermost) addEventMetadata(event Event) Event {
	var user, actor *model.User
	var channel *model.Channel
	if event.UserName == "" && event.UserID != "" {
		user, _ = m.api.GetUser(event.UserID)
	} else {
		user, _ = m.api.GetUserByUsername(event.UserName)
	}

	if user != nil {
		event.UserID = user.Id
		event.UserName = user.Username
		event.UserRole = user.Roles
	}

	//member, _ := m.api.GetTeamMember(m.team.Id, user.Id)
	//event.UserRole += " " + member.Roles

	if event.ActorName == "" && event.ActorID != "" {
		actor, _ = m.api.GetUser(event.ActorID)
	} else if event.ActorName != "" && event.ActorID == "" {
		actor, _ = m.api.GetUserByUsername(event.ActorName)
	}

	if actor != nil {
		event.ActorID = actor.Id
		event.ActorName = actor.Username
		event.ActorRole = actor.Roles

		//member, _ = m.api.GetTeamMember(m.team.Id, actor.Id)
		//event.ActorRole += " " + member.Roles
	}

	if event.ChannelID != "" {
		channel, _ = m.api.GetChannel(event.ChannelID)
	} else {
		//channel, _ = m.api.GetChannelByName(event.ChannelName, m.team.Id)
	}

	if channel != nil {
		event.ChannelID = channel.Id
		event.ChannelName = channel.Name
		event.ChannelType = channel.Type
	}

	return event
}

func (m *Mattermost) GetEmailByUsername(username string) string {
	userInfo, resp := m.api.GetUserByUsername(username)
	log.Println("[+] eMail for Mattermost LookUp:", userInfo.Email)
	if resp != nil {
		log.Println("[!] Unable to get UserInfo:", resp)
		return ""
	}

	return userInfo.Email
}

func (m *Mattermost) MessagePublic(channelID, message string) bool {
	post := &model.Post{
		ChannelId: channelID,
		Message:   message,
		UserId:    m.botUser.Id,
		Props: model.StringInterface{
			"from_claptrap": true,
		},
	}
	m.api.CreatePost(post)

	return true
}

func (m *Mattermost) MessageUser(userID, message string) bool {
	channel, _ := m.api.GetDirectChannel(userID, m.botUser.Id)
	m.MessagePublic(channel.Id, message)
	return true
}

func (m *Mattermost) MessageEphemeral(userID, channelID, message string) bool {
	post := &model.Post{
		ChannelId: channelID,
		Message:   message,
	}
	m.api.SendEphemeralPost(userID, post)

	return true
}

func (m *Mattermost) InviteUser(userID, channelID string) bool {
	m.api.AddChannelMember(channelID, userID)
	return true
}

func (m *Mattermost) KickUser(userID, channelID string) bool {
	m.api.DeleteChannelMember(channelID, userID)
	return true
}

func (m *Mattermost) DeleteMessage(event Event) bool {
	m.api.DeletePost(event.PostID)
	return true
}

func (m *Mattermost) ReplaceMessagePlaceholders(event Event, message string) string {
	message = strings.Replace(message, "{User_Name}", event.UserName, 1)
	message = strings.Replace(message, "{Actor_Name}", event.ActorName, 1)
	message = strings.Replace(message, "{Bot_Name}", m.botUser.Username, 1)
	message = strings.Replace(message, "{Channel_Name}", event.ChannelName, 1)
	return message
}
