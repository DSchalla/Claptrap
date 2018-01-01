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
	log.Printf("[+] Executing 'MessageChannelResponse'")
	message := h.ReplaceMessagePlaceholders(event, k.Message)
	return h.MessagePublic(event.ChannelID, message)
}

type MessageUserResponse struct {
	UserID  string
	Message string
}

func (k MessageUserResponse) Execute(h ResponseHandler, event Event) bool {
	log.Printf("[+] Executing 'MessageUserResponse'")

	userID := ""

	if k.UserID == "" {
		userID = event.UserID
	} else {
		userID = k.UserID
	}
	message := h.ReplaceMessagePlaceholders(event, k.Message)
	return h.MessageUser(userID, message)
}

type InviteUserResponse struct {
	ChannelID string
	UserID    string
}

func (k InviteUserResponse) Execute(h ResponseHandler, event Event) bool {
	log.Printf("[+] Executing 'InviteUserResponse'")

	userID := ""

	if k.UserID == "" {
		userID = event.UserID
	} else {
		userID = k.UserID
	}

	return h.InviteUser(event.ChannelID, userID)
}

type KickUserResponse struct {
	ChannelID string
	UserID    string
}

func (k KickUserResponse) Execute(h ResponseHandler, event Event) bool {
	log.Printf("[+] Executing 'KickUserResponse'")

	userID := ""

	if k.UserID == "" {
		userID = event.UserID
	} else {
		userID = k.UserID
	}

	return h.KickUser(userID, event.ChannelID)
}

type DeleteMessageResponse struct {
}

func (d DeleteMessageResponse) Execute(h ResponseHandler, event Event) bool {
	log.Printf("[+] Executing 'DeleteMessage'")
	return h.DeleteMessage(event.ChannelID, event.Timestamp)
}
