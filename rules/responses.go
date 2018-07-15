package rules

import (
	"github.com/DSchalla/Claptrap/provider"
	"log"
)

type Response interface {
	GetName() string
	Execute(p provider.Provider, event provider.Event) bool
}

func NewMessageChannelResponse(channelID, message string) (*MessageChannelResponse, error) {
	return &MessageChannelResponse{channelID, message}, nil
}

type MessageChannelResponse struct {
	ChannelID string
	Message   string
}

func (m MessageChannelResponse) GetName() string{
	return "MessageChannelResponse"
}

func (m MessageChannelResponse) Execute(p provider.Provider, event provider.Event) bool {
	message := p.ReplaceMessagePlaceholders(event, m.Message)
	channelID := ""
	if len(m.ChannelID) == 0 {
		channelID = event.ChannelID
	} else {
		channelID = m.ChannelID
	}

	log.Printf("[+] Executing 'MessageChannelResponse' | ChannelID: %s \n", channelID)
	return p.MessagePublic(channelID, message)
}

func NewMessageUserResponse(userID, message string) (*MessageUserResponse, error) {
	return &MessageUserResponse{userID, message}, nil
}

type MessageUserResponse struct {
	UserID  string
	Message string
}

func (m MessageUserResponse) GetName() string{
	return "MessageUserResponse"
}

func (m MessageUserResponse) Execute(p provider.Provider, event provider.Event) bool {
	userID := ""

	if m.UserID == "" {
		userID = event.UserID
	} else {
		userID = m.UserID
	}
	message := p.ReplaceMessagePlaceholders(event, m.Message)

	log.Printf("[+] Executing 'MessageUserResponse' | UserID: %s \n", userID)

	return p.MessageUser(userID, message)
}

func NewMessageEphemeralResponse(message string) (*MessageEphemeralResponse, error) {
	return &MessageEphemeralResponse{message}, nil
}

type MessageEphemeralResponse struct {
	Message string
}

func (m MessageEphemeralResponse) GetName() string{
	return "MessageEphemeralResponse"
}

func (m MessageEphemeralResponse) Execute(p provider.Provider, event provider.Event) bool {

	userID := event.UserID
	message := p.ReplaceMessagePlaceholders(event, m.Message)

	log.Printf("[+] Executing 'MessageEphemeralResponse' | UserID: %s \n", userID)

	return p.MessageEphemeral(userID, event.ChannelID, message)
}

func NewInviteUserResponse(channelID, userID string) (*InviteUserResponse, error) {
	return &InviteUserResponse{channelID, userID}, nil
}

type InviteUserResponse struct {
	ChannelID string
	UserID    string
}

func (i InviteUserResponse) GetName() string{
	return "InviteUserResponse"
}

func (i InviteUserResponse) Execute(p provider.Provider, event provider.Event) bool {

	userID := ""

	if i.UserID == "" {
		userID = event.UserID
	} else {
		userID = i.UserID
	}

	log.Printf("[+] Executing 'InviteUserResponse' | ChannelID: %s | UserID: %s\n", i.ChannelID, userID)
	return p.InviteUser(userID, i.ChannelID)
}

func NewKickUserResponse(channelID, userID string) (*KickUserResponse, error) {
	return &KickUserResponse{channelID, userID}, nil
}

type KickUserResponse struct {
	ChannelID string
	UserID    string
}

func (k KickUserResponse) GetName() string{
	return "KickUserResponse"
}

func (k KickUserResponse) Execute(p provider.Provider, event provider.Event) bool {
	userID := ""

	if k.UserID == "" {
		userID = event.UserID
	} else {
		userID = k.UserID
	}

	channelID := ""
	channelName := ""
	if len(k.ChannelID) == 0 {
		channelID = event.ChannelID
		channelName = event.ChannelName
	} else {
		channelID = k.ChannelID
		channelName = "?"
	}

	log.Printf("[+] Executing 'KickUserResponse' | Channel: %s (%s) | UserID: %s\n", channelName, channelID, userID)
	return p.KickUser(userID, channelID)
}

func NewDeleteMessageResponse() (*DeleteMessageResponse, error) {
	return &DeleteMessageResponse{}, nil
}

type DeleteMessageResponse struct {
}

func (d DeleteMessageResponse) GetName() string{
	return "DeleteMessageResponse"
}

func (d DeleteMessageResponse) Execute(p provider.Provider, event provider.Event) bool {
	log.Printf("[+] Executing 'DeleteMessage' | Channel: %s (%s) | PostID: %s\n", event.ChannelName, event.ChannelID, event.PostID)
	return p.DeleteMessage(event)
}
