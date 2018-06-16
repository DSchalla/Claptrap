package provider

import (
	"github.com/mattermost/mattermost-server/model"
	"log"
	"os"
	"strings"
)

type Mattermost struct {
	client    *model.Client4
	socket    *model.WebSocketClient
	apiUrl    string
	team      *model.Team
	botUser   *model.User
	eventChan chan Event
}

func NewMattermost(apiUrl, username, password, team string) *Mattermost {
	m := Mattermost{}
	m.apiUrl = apiUrl
	m.client = model.NewAPIv4Client("http://" + apiUrl)
	if user, resp := m.client.Login(username, password); resp.Error != nil {
		log.Println("[!] Error authenticating against Mattermost API")
		os.Exit(1)
	} else {
		m.botUser = user
	}
	m.team, _ = m.client.GetTeamByName(team, "")
	m.eventChan = make(chan Event, 10)
	return &m
}

func (m *Mattermost) Connect() bool {
	var err *model.AppError
	m.socket, err = model.NewWebSocketClient4("ws://"+m.apiUrl, m.client.AuthToken)
	if err != nil {
		log.Printf("[!] Error connecting to the Mattermost WS: %s\n", err.Message)
		return false
	}
	m.socket.Listen()
	log.Println("[+] Mattermost Websocket connection established")

	return true
}

func (m *Mattermost) Reconnect() bool {
	return m.Connect()
}

func (m *Mattermost) ListenForEvents() {
	var unifiedEvent Event
	for msg := range m.socket.EventChannel {
		unifiedEvent = Event{}
		switch msg.Event {
		case "posted":
			unifiedEvent = m.handleMessageEvent(msg)
		case "user_removed":
			unifiedEvent = m.handleUserRemovedEvent(msg)
		default:
			unifiedEvent = Event{}
		}

		if unifiedEvent.UserID == m.botUser.Id {
			continue
		}

		if unifiedEvent.Type != "" {
			m.eventChan <- unifiedEvent
		}
	}
}

func (m *Mattermost) handleMessageEvent(event *model.WebSocketEvent) Event {

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post.Type == "system_add_to_channel" {
		return m.handleUserInviteEvent(event)
	} else if post.Type == "system_join_channel" {
		return m.handleUserJoinEvent(event)
	}

	unifiedEvent := Event{}
	unifiedEvent.Type = "message"
	unifiedEvent.UserName = event.Data["sender_name"].(string)
	unifiedEvent.ChannelName = event.Data["channel_name"].(string)
	unifiedEvent.PostID = post.Id
	unifiedEvent.UserID = post.UserId
	unifiedEvent.ChannelID = post.ChannelId
	unifiedEvent.Timestamp = post.CreateAt
	unifiedEvent.Text = post.Message
	unifiedEvent = m.addEventMetadata(unifiedEvent)
	return unifiedEvent
}

func (m *Mattermost) handleUserInviteEvent(event *model.WebSocketEvent) Event {
	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	unifiedEvent := Event{}
	unifiedEvent.Type = "user_add"
	unifiedEvent.PostID = post.Id
	unifiedEvent.ChannelID = post.ChannelId
	unifiedEvent.UserName = post.Props["addedUsername"].(string)
	unifiedEvent.ActorName = post.Props["username"].(string)
	unifiedEvent = m.addEventMetadata(unifiedEvent)
	return unifiedEvent
}

func (m *Mattermost) handleUserJoinEvent(event *model.WebSocketEvent) Event {
	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	unifiedEvent := Event{}
	unifiedEvent.Type = "user_add"
	unifiedEvent.PostID = post.Id
	unifiedEvent.ChannelID = post.ChannelId
	unifiedEvent.UserName = post.Props["username"].(string)
	unifiedEvent = m.addEventMetadata(unifiedEvent)
	return unifiedEvent
}

func (m *Mattermost) handleUserRemovedEvent(event *model.WebSocketEvent) Event {
	unifiedEvent := Event{}
	unifiedEvent.Type = "user_remove"
	unifiedEvent.UserID = event.Data["user_id"].(string)
	unifiedEvent.ActorID = event.Data["remover_id"].(string)
	unifiedEvent.ChannelID = event.Broadcast.ChannelId
	unifiedEvent = m.addEventMetadata(unifiedEvent)
	return unifiedEvent
}

func (m *Mattermost) addEventMetadata(event Event) Event {
	var user, actor *model.User
	var channel *model.Channel
	if event.UserName == "" && event.UserID != "" {
		user, _ = m.client.GetUser(event.UserID, "")
	} else {
		user, _ = m.client.GetUserByUsername(event.UserName, "")
	}

	if user != nil {
		event.UserID = user.Id
		event.UserName = user.Username
		event.UserRole = user.Roles
	}

	member, _ := m.client.GetTeamMember(m.team.Id, user.Id, "")
	event.UserRole += " " + member.Roles

	if event.ActorName == "" && event.ActorID != "" {
		actor, _ = m.client.GetUser(event.ActorID, "")
	} else if event.ActorName != "" && event.ActorID == "" {
		actor, _ = m.client.GetUserByUsername(event.ActorName, "")
	}

	if actor != nil {
		event.ActorID = actor.Id
		event.ActorName = actor.Username
		event.ActorRole = actor.Roles

		member, _ = m.client.GetTeamMember(m.team.Id, actor.Id, "")
		event.ActorRole += " " + member.Roles
	}

	if event.ChannelID != "" {
		channel, _ = m.client.GetChannel(event.ChannelID, "")
	} else {
		channel, _ = m.client.GetChannelByName(event.ChannelName, m.team.Id, "")
	}

	if channel != nil {
		event.ChannelID = channel.Id
		event.ChannelName = channel.Name
		event.ChannelType = channel.Type
	}

	return event
}

/* ToDo: Implement IsAlive handler */
func (m *Mattermost) IsAlive() bool {
	return true
}

func (m *Mattermost) GetEvents() <-chan Event {
	return m.eventChan
}

func (m *Mattermost) AutoJoinAllChannel() bool {
	channels, _ := m.client.GetPublicChannelsForTeam(m.team.Id, 0, 500, "")
	log.Println("[+] Triggering Autojoin of Channels")
	log.Printf("[+] Getting Channel List -> Found %d channels\n", len(channels))
	for _, channel := range channels {
		m.client.AddChannelMember(channel.Id, m.botUser.Id)
		log.Printf("[+] Joining Channel '%s'\n", channel.Name)
	}
	return true
}

func (m *Mattermost) GetEmailByUsername(username string) string {
	userInfo, resp := m.client.GetUserByUsername(username, "")
	log.Println("[+] eMail for Mattermost LookUp:", userInfo.Email)
	if resp.Error != nil {
		log.Println("[!] Unable to get UserInfo:", resp.Error)
		return ""
	}

	return userInfo.Email
}

func (m *Mattermost) MessagePublic(channelID, message string) bool {
	post := &model.Post{
		ChannelId: channelID,
		Message:   message,
	}
	m.client.CreatePost(post)

	return true
}

func (m *Mattermost) MessageUser(userID, message string) bool {
	channel, _ := m.client.CreateDirectChannel(userID, m.botUser.Id)
	m.MessagePublic(channel.Id, message)
	return true
}

func (m *Mattermost) InviteUser(userID, channelID string) bool {
	m.client.AddChannelMember(channelID, userID)
	return true
}

func (m *Mattermost) KickUser(userID, channelID string) bool {
	m.client.RemoveUserFromChannel(channelID, userID)
	return true
}

func (m *Mattermost) DeleteMessage(event Event) bool {
	m.client.DeletePost(event.PostID)
	return true
}

func (m *Mattermost) ReplaceMessagePlaceholders(event Event, message string) string {
	message = strings.Replace(message, "{User_Name}", event.UserName, 1)
	message = strings.Replace(message, "{Actor_Name}", event.ActorName, 1)
	message = strings.Replace(message, "{Bot_Name}", m.botUser.Username, 1)
	message = strings.Replace(message, "{Channel_Name}", event.ChannelName, 1)
	return message
}
