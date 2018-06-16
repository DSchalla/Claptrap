package rules

import (
	"github.com/DSchalla/Claptrap/provider"
	"log"
)

type Response interface {
	Execute(p provider.Provider, event provider.Event) bool
}

func NewMessageChannelResponse(channelID, message string) (*MessageChannelResponse, error) {
	return &MessageChannelResponse{channelID, message}, nil
}

type MessageChannelResponse struct {
	channelID string
	message   string
}

func (k MessageChannelResponse) Execute(p provider.Provider, event provider.Event) bool {
	message := p.ReplaceMessagePlaceholders(event, k.message)
	channelID := ""
	if len(k.channelID) == 0 {
		channelID = event.ChannelID
	} else {
		channelID = k.channelID
	}

	log.Printf("[+] Executing 'MessageChannelResponse' | ChannelID: %s \n", channelID)
	return p.MessagePublic(channelID, message)
}

func NewMessageUserResponse(userID, message string) (*MessageUserResponse, error) {
	return &MessageUserResponse{userID, message}, nil
}

type MessageUserResponse struct {
	userID  string
	message string
}

func (k MessageUserResponse) Execute(p provider.Provider, event provider.Event) bool {
	userID := ""

	if k.userID == "" {
		userID = event.UserID
	} else {
		userID = k.userID
	}
	message := p.ReplaceMessagePlaceholders(event, k.message)

	log.Printf("[+] Executing 'MessageUserResponse' | UserID: %s \n", userID)

	return p.MessageUser(userID, message)
}

func NewInviteUserResponse(channelID, userID string) (*InviteUserResponse, error) {
	return &InviteUserResponse{channelID, userID}, nil
}

type InviteUserResponse struct {
	channelID string
	userID    string
}

func (k InviteUserResponse) Execute(p provider.Provider, event provider.Event) bool {

	userID := ""

	if k.userID == "" {
		userID = event.UserID
	} else {
		userID = k.userID
	}

	log.Printf("[+] Executing 'InviteUserResponse' | ChannelID: %s | UserID: %s\n", k.channelID, userID)
	return p.InviteUser(userID, k.channelID)
}

func NewKickUserResponse(channelID, userID string) (*KickUserResponse, error) {
	return &KickUserResponse{channelID, userID}, nil
}

type KickUserResponse struct {
	ChannelID string
	UserID    string
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

func (d DeleteMessageResponse) Execute(p provider.Provider, event provider.Event) bool {
	log.Printf("[+] Executing 'DeleteMessage' | Channel: %s (%s) | PostID: %s\n", event.ChannelName, event.ChannelID, event.PostID)
	return p.DeleteMessage(event)
}
