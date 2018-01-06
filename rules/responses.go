package rules

import (
	"log"
)

type ResponseHandler interface {
	MessagePublic(channelID, message string) bool
	MessageUser(userID, message string) bool
	InviteUser(userID, channelID string) bool
	KickUser(userID, channelID string) bool
	DeleteMessage(channelID, timestamp string) bool
	ReplaceMessagePlaceholders(event Event, message string) string
}

type Response interface {
	Execute(h ResponseHandler, event Event) bool
}

type MessageChannelResponse struct {
	ChannelID string
	Message   string
}

func (k MessageChannelResponse) Execute(h ResponseHandler, event Event) bool {
	message := h.ReplaceMessagePlaceholders(event, k.Message)
	channelID := ""
	if len(k.ChannelID) == 0{
		channelID = event.ChannelID
	} else {
		channelID = k.ChannelID
	}

	log.Printf("[+] Executing 'MessageChannelResponse' | ChannelID: %s \n", channelID)
	return h.MessagePublic(channelID, message)
}

type MessageUserResponse struct {
	UserID  string
	Message string
}

func (k MessageUserResponse) Execute(h ResponseHandler, event Event) bool {
	userID := ""

	if k.UserID == "" {
		userID = event.UserID
	} else {
		userID = k.UserID
	}
	message := h.ReplaceMessagePlaceholders(event, k.Message)

	log.Printf("[+] Executing 'MessageUserResponse' | UserID: %s \n", userID)

	return h.MessageUser(userID, message)
}

type InviteUserResponse struct {
	ChannelID string
	UserID    string
}

func (k InviteUserResponse) Execute(h ResponseHandler, event Event) bool {

	userID := ""

	if k.UserID == "" {
		userID = event.UserID
	} else {
		userID = k.UserID
	}

	log.Printf("[+] Executing 'InviteUserResponse' | ChannelID: %s | UserID: %s", k.ChannelID, userID)
	return h.InviteUser(userID, k.ChannelID)
}

type KickUserResponse struct {
	ChannelID string
	UserID    string
}

func (k KickUserResponse) Execute(h ResponseHandler, event Event) bool {
	userID := ""

	if k.UserID == "" {
		userID = event.UserID
	} else {
		userID = k.UserID
	}

	channelID := ""
	if len(k.ChannelID) == 0{
		channelID = event.ChannelID
	} else {
		channelID = k.ChannelID
	}

	log.Printf("[+] Executing 'KickUserResponse' | ChannelID: %s | UserID: %s", channelID, userID)
	return h.KickUser(userID, channelID)
}

type DeleteMessageResponse struct {
}

func (d DeleteMessageResponse) Execute(h ResponseHandler, event Event) bool {
	log.Printf("[+] Executing 'DeleteMessage' | ChannelID: %s | Timestamp: %s", event.ChannelID, event.Timestamp)
	return h.DeleteMessage(event.ChannelID, event.Timestamp)
}
