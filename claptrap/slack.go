package claptrap

import (
	"github.com/nlopes/slack"
	"log"
	"strings"
)

func NewSlackHandler(botToken, adminToken string) *SlackHandler {
	slackHandler := SlackHandler{}
	slackHandler.BotAPI = slack.New(botToken)
	slackHandler.AdminAPI = slack.New(adminToken)
	return &slackHandler
}

type SlackHandler struct {
	BotAPI   *slack.Client
	AdminAPI *slack.Client
}

func (h *SlackHandler) StartRTM() *slack.RTM {
	rtm := h.BotAPI.NewRTM()
	go rtm.ManageConnection()
	log.Println("[+] Slack RTM connection established")
	return rtm
}

func NewSlackResponseHandler(rtm *slack.RTM, adminAPI *slack.Client) *SlackResponseHandler {
	handler := SlackResponseHandler{
		botRTM:   rtm,
		adminAPI: adminAPI,
	}
	return &handler
}

type SlackResponseHandler struct {
	botRTM   *slack.RTM
	adminAPI *slack.Client
}

func (s SlackResponseHandler) MessagePublic(channelID, message string) bool {
	params := slack.PostMessageParameters{
		AsUser: true,
	}
	_, _, err := s.botRTM.PostMessage(channelID, message, params)

	if err != nil {
		return false
	}

	return true
}

func (s SlackResponseHandler) MessageUser(userID, message string) bool {
	return s.MessagePublic(userID, message)
}

func (s SlackResponseHandler) InviteUser(userID, channelID string) bool {
	var err error

	if strings.HasPrefix(channelID, "G") {
		_, _, err = s.adminAPI.InviteUserToGroup(channelID, userID)
	} else {
		_, err = s.adminAPI.InviteUserToChannel(channelID, userID)
	}

	if err != nil {
		return false
	}

	return true
}

func (s SlackResponseHandler) KickUser(userID, channelID string) bool {
	var err error

	if strings.HasPrefix(channelID, "G") {
		err = s.adminAPI.KickUserFromGroup(channelID, userID)
	} else {
		err = s.adminAPI.KickUserFromChannel(channelID, userID)
	}

	if err != nil {
		return false
	}

	return true
}

func (s SlackResponseHandler) DeleteMessage(channelID, timestamp string) bool {
	_, _, err := s.adminAPI.DeleteMessage(channelID, timestamp)

	if err != nil {
		return false
	}

	return true
}
