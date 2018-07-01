package provider

import (
	"strings"
	"github.com/mattermost/mattermost-server/model"
)

type messagePublicLog struct {
	ChannelID string
	Message   string
}

type messageUserLog struct {
	UserID  string
	Message string
}

type inviteUserLog struct {
	ChannelID string
	UserID    string
}

type kickUserLog struct {
	ChannelID string
	UserID    string
}

type deleteMessageLog struct {
	PostID string
}

type Debug struct {
	MessagePublicLog []messagePublicLog
	MessageUserLog   []messageUserLog
	InviteUserLog    []inviteUserLog
	KickUserLog      []kickUserLog
	DeleteMessageLog []deleteMessageLog
}

func NewDebug() *Debug {
	d := Debug{}
	return &d
}

func (Debug) AutoJoinAllChannel() error {
	panic("implement me")
}

func (Debug) GetEmailByUsername(username string) string {
	panic("implement me")
}

func (Debug) NormalizeMessageEvent(post *model.Post) Event {
	panic("implement me")
}

func (Debug) NormalizeUserJoinEvent(post *model.Post) Event {
	panic("implement me")
}

func (Debug) NormalizeUserInviteEvent(post *model.Post) Event {
	panic("implement me")
}

func (d *Debug) MessagePublic(channelID, message string) bool {
	d.MessagePublicLog = append(d.MessagePublicLog, messagePublicLog{channelID, message})
	return true
}

func (d *Debug) MessageUser(userID, message string) bool {
	d.MessageUserLog = append(d.MessageUserLog, messageUserLog{userID, message})
	return true
}

func (Debug) MessageEphemeral(userID, channelID, message string) bool {
	panic("implement me")
}

func (d *Debug) InviteUser(userID, channelID string) bool {
	d.InviteUserLog = append(d.InviteUserLog, inviteUserLog{channelID, userID})
	return true
}

func (d *Debug) KickUser(userID, channelID string) bool {
	d.KickUserLog = append(d.KickUserLog, kickUserLog{channelID, userID})
	return true
}

func (d *Debug) DeleteMessage(event Event) bool {
	d.DeleteMessageLog = append(d.DeleteMessageLog, deleteMessageLog{event.PostID})
	return true
}

func (Debug) ReplaceMessagePlaceholders(event Event, message string) string {
	message = strings.Replace(message, "{User_Name}", event.UserName, 1)
	message = strings.Replace(message, "{Actor_Name}", event.ActorName, 1)
	message = strings.Replace(message, "{Bot_Name}", "Bot1", 1)
	message = strings.Replace(message, "{Channel_Name}", event.ChannelName, 1)
	return message
}
