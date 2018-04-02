package provider

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

func (h *SlackHandler) AutoJoinAllChannel(botID string) {
	channels, err := h.AdminAPI.GetChannels(true)

	if err != nil {
		log.Println("[!] Error getting channel list: ", err)
	}

	for _, channel := range channels {

		if channel.IsMember {
			continue
		}

		log.Printf("[+] Attempting to Join Channel '%s' (%s) as Admin\n", channel.Name, channel.ID)
		_, err := h.AdminAPI.JoinChannel(channel.Name)

		if err != nil {
			log.Println("[!] Error joining channel as admin:", err)
		}

		log.Printf("[+] Attempting to Invite Bot to Channel '%s' (%s)\n", channel.Name, channel.ID)
		_, err = h.AdminAPI.InviteUserToChannel(channel.ID, botID)

		if err != nil {
			log.Println("[!] Error inviting bot to channel as admin:", err)
		}
	}
}

func NewSlackResponseHandler(rtm *slack.RTM, adminAPI *slack.Client) *SlackResponseHandler {
	handler := SlackResponseHandler{
		BotRTM:   rtm,
		AdminAPI: adminAPI,
	}
	return &handler
}

type SlackResponseHandler struct {
	BotRTM   *slack.RTM
	AdminAPI *slack.Client
}

func (s SlackResponseHandler) MessagePublic(channelID, message string) bool {
	params := slack.PostMessageParameters{
		AsUser: true,
	}
	_, _, err := s.BotRTM.PostMessage(channelID, message, params)

	if err != nil {
		log.Println("[!] Response from API Endpoint:", err)
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
		_, _, err = s.AdminAPI.InviteUserToGroup(channelID, userID)
	} else {
		_, err = s.AdminAPI.InviteUserToChannel(channelID, userID)
	}

	if err != nil {
		log.Println("[!] Response from API Endpoint:", err)
		return false
	}

	return true
}

func (s SlackResponseHandler) KickUser(userID, channelID string) bool {
	var err error

	if strings.HasPrefix(channelID, "G") {
		err = s.AdminAPI.KickUserFromGroup(channelID, userID)
	} else {
		err = s.AdminAPI.KickUserFromChannel(channelID, userID)
	}

	if err != nil {
		log.Println("[!] Response from API Endpoint:", err)
		return false
	}

	return true
}

func (s SlackResponseHandler) DeleteMessage(channelID, timestamp string) bool {
	_, _, err := s.AdminAPI.DeleteMessage(channelID, timestamp)

	if err != nil {
		log.Println("[!] Response from API Endpoint:", err)
		return false
	}

	return true
}

func (s SlackResponseHandler) ReplaceMessagePlaceholders(event Event, message string) string {
	botInfo := s.BotRTM.GetInfo()
	message = strings.Replace(message, "{Sender_Name}", "<@"+event.UserID+">", 1)
	message = strings.Replace(message, "{Bot_Name}", "<@"+botInfo.User.ID+">", 1)
	message = strings.Replace(message, "{Channel_Name}", event.ChannelName, 1)
	return message
}
