package provider

import "github.com/mattermost/mattermost-server/model"

type Provider interface {
	AutoJoinAllChannel() error

	// Normalization
	NormalizeMessageEvent(post *model.Post) Event
	NormalizeTeamJoinEvent(teamMember *model.TeamMember, actor *model.User) Event
	NormalizeChannelJoinEvent(channelMember *model.ChannelMember, actor *model.User) Event
	NormalizeChannelLeaveEvent(channelMember *model.ChannelMember, actor *model.User) Event

	// Get Information
	GetEmailByUsername(username string) string

	// Actions
	MessagePublic(channelID, message string) bool
	MessageUser(userID, message string) bool
	MessageEphemeral(userID, channelID, message string) bool
	InviteUser(userID, channelID string) bool
	KickUser(userID, channelID string) bool
	DeleteMessage(event Event) bool
	ReplaceMessagePlaceholders(event Event, message string) string
}
