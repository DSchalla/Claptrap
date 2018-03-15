package claptrap

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/DSchalla/Claptrap/rules"
	"log"
	"os"
	"fmt"
	"strings"
)

type MattermostHandler struct {
	Client  *model.Client4
	Socket  *model.WebSocketClient
	apiUrl  string
	Team    *model.Team
	BotUser *model.User
}

func NewMattermostHandler (apiUrl, username, password, team string) *MattermostHandler{
	mh := MattermostHandler{
	}
	mh.apiUrl = apiUrl
	mh.Client = model.NewAPIv4Client("http://" + apiUrl)
	if user, resp := mh.Client.Login(username, password); resp.Error != nil {
		fmt.Println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		os.Exit(1)
	} else {
		mh.BotUser = user
	}
	mh.Team, _ = mh.Client.GetTeamByName(team, "")
	return &mh
}

func (m *MattermostHandler) StartWS() {
	var err *model.AppError
	m.Socket, err = model.NewWebSocketClient4("ws://" + m.apiUrl, m.Client.AuthToken)
	if err != nil {
		log.Printf("[!] Error connecting to the Mattermost WS: %s\n", err.Message)
		return
	}
	m.Socket.Listen()
	log.Println("[+] Mattermost Websocket connection established")
}

func (m *MattermostHandler) AutoJoinAllChannel() {
	channels, _ := m.Client.GetPublicChannelsForTeam(m.Team.Id, 0, 500, "")
	log.Println("[+] Triggering Autojoin of Channels")
	log.Printf("[+] Getting Channel List -> Found %d channels\n", len(channels))
	for _, channel := range channels {
		m.Client.AddChannelMember(channel.Id, m.BotUser.Id)
		log.Printf("[+] Joining Channel '%s'\n", channel.Name)
	}
}

func (m *MattermostHandler) GetMessages() <-chan *model.WebSocketEvent {
	return m.Socket.EventChannel
}

type MattermostResponseHandler struct {
	client *model.Client4
	botUser *model.User
}

func NewMattermostResponseHandler(client *model.Client4, botUser *model.User) *MattermostResponseHandler {
	handler := MattermostResponseHandler{
		client: client,
		botUser: botUser,
	}
	return &handler
}

func (m MattermostResponseHandler) MessagePublic(channelID, message string) bool {
	post := &model.Post{
		ChannelId: channelID,
		Message: message,
	}
	m.client.CreatePost(post)

	return true
}

func (m MattermostResponseHandler) MessageUser(userID, message string) bool {
	channel,_ := m.client.CreateDirectChannel(userID, m.botUser.Id)
	m.MessagePublic(channel.Id, message)
	return true
}

func (m MattermostResponseHandler) InviteUser(userID, channelID string) bool {
	m.client.AddChannelMember(channelID, userID)
	return true
}

func (m MattermostResponseHandler) KickUser(userID, channelID string) bool {
	m.client.RemoveUserFromChannel(channelID, userID)
	return true
}

func (m MattermostResponseHandler) DeleteMessage(postID string) bool {
	m.client.DeletePost(postID)
	return true
}

func (m MattermostResponseHandler) ReplaceMessagePlaceholders(event rules.Event, message string) string {
	message = strings.Replace(message, "{User_Name}", event.UserName, 1)
	message = strings.Replace(message, "{Actor_Name}", event.ActorName, 1)
	message = strings.Replace(message, "{Bot_Name}", m.botUser.Username, 1)
	message = strings.Replace(message, "{Channel_Name}", event.ChannelName, 1)
	return message
}
