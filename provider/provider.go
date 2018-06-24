package provider

import "github.com/mattermost/mattermost-server/model"

type Provider interface {
	AutoJoinAllChannel() bool

	// Normalization
	NormalizeMessageEvent(post *model.Post) Event
	NormalizeUserInviteEvent(post *model.Post) Event
	NormalizeUserJoinEvent(post *model.Post) Event

	// Get Information
	GetEmailByUsername(username string) string

	// Actions
	MessagePublic(channelID, message string) bool
	MessageUser(userID, message string) bool
	InviteUser(userID, channelID string) bool
	KickUser(userID, channelID string) bool
	DeleteMessage(event Event) bool
	ReplaceMessagePlaceholders(event Event, message string) string
}
