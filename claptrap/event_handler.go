package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"log"
	"github.com/mattermost/mattermost-server/model"
	"fmt"
	"strings"
)

type EventHandler struct {
	mhHandler	*MattermostHandler
	ruleEngine *rules.Engine
}

func NewEventHandler(mhHandler *MattermostHandler, ruleEngine *rules.Engine) *EventHandler {
	eh := EventHandler{
		mhHandler: mhHandler,
	}
	eh.ruleEngine = ruleEngine
	return &eh
}

func (eh *EventHandler) Start() {
	log.Println("[+] Event Handler started and listening")
	var unifiedEvent rules.Event
	for msg := range eh.mhHandler.GetMessages() {
		switch msg.Event {
			case "posted": unifiedEvent = eh.handleMessageEvent(msg)
			case "user_removed": unifiedEvent = eh.handleUserRemovedEvent(msg)
		}

		if unifiedEvent.UserID == eh.mhHandler.BotUser.Id {
			continue
		}

		if unifiedEvent.Type != "" {
			eh.ruleEngine.EvaluateEvent(unifiedEvent)
		}

		fmt.Printf("Unexpected: %+v\n", msg)
		fmt.Printf("%+v\n", msg.Broadcast)
	}
}

func (eh *EventHandler) handleMessageEvent(event *model.WebSocketEvent) rules.Event{

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post.Type == "system_add_to_channel"{
		return eh.handleUserAddEvent(event)
	}
	unifiedEvent := rules.Event{}
	unifiedEvent.Type = "message"
	unifiedEvent.UserName = event.Data["sender_name"].(string)
	unifiedEvent.ChannelName = event.Data["channel_name"].(string)
	unifiedEvent.PostID = post.Id
	unifiedEvent.UserID = post.UserId
	unifiedEvent.ChannelID = post.ChannelId
	unifiedEvent.Timestamp = post.CreateAt
	unifiedEvent.Text = post.Message
	unifiedEvent = eh.addEventMetadata(unifiedEvent)
	return unifiedEvent
}

func (eh *EventHandler) handleUserAddEvent(event *model.WebSocketEvent) rules.Event{
	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	unifiedEvent := rules.Event{}
	unifiedEvent.Type = "user_add"
	unifiedEvent.PostID = post.Id
	unifiedEvent.ChannelID = post.ChannelId
	unifiedEvent.UserName = post.Props["addedUsername"].(string)
	unifiedEvent.ActorName = post.Props["username"].(string)
	unifiedEvent = eh.addEventMetadata(unifiedEvent)
	return unifiedEvent
}

func (eh *EventHandler) handleUserRemovedEvent(event *model.WebSocketEvent) rules.Event{
	unifiedEvent := rules.Event{}
	unifiedEvent.Type = "user_remove"
	unifiedEvent.UserID = event.Data["user_id"].(string)
	unifiedEvent.ActorID = event.Data["remover_id"].(string)
	unifiedEvent.ChannelID = event.Broadcast.ChannelId
	unifiedEvent = eh.addEventMetadata(unifiedEvent)
	return unifiedEvent
}

func (eh *EventHandler) addEventMetadata(event rules.Event) rules.Event {
	var user, actor *model.User
	var channel *model.Channel
	if event.UserName == "" && event.UserID != "" {
		user, _ = eh.mhHandler.Client.GetUser(event.UserID, "")
	} else  {
		user, _ = eh.mhHandler.Client.GetUserByUsername(event.UserName, "")
	}

	if user != nil {
		event.UserID = user.Id
		event.UserName = user.Username
		event.UserRole = user.Roles
	}

	member, _ := eh.mhHandler.Client.GetTeamMember(eh.mhHandler.Team.Id, user.Id, "")
	event.UserRole += " " + member.Roles

	if event.ActorName == "" && event.ActorID != "" {
		actor, _ = eh.mhHandler.Client.GetUser(event.ActorID, "")
	} else if event.ActorName != "" && event.ActorID == "" {
		actor, _ = eh.mhHandler.Client.GetUserByUsername(event.ActorName, "")
	}

	if actor != nil {
		event.ActorID = actor.Id
		event.ActorName = actor.Username
		event.ActorRole = actor.Roles

		member, _ = eh.mhHandler.Client.GetTeamMember(eh.mhHandler.Team.Id, actor.Id, "")
		event.ActorRole += " " + member.Roles
	}

	if event.ChannelName == "" && event.ChannelID != "" {
		channel, _ = eh.mhHandler.Client.GetChannel(event.ChannelID, "")
	} else if event.ChannelName != "" && event.ChannelID == "" {
		channel, _ = eh.mhHandler.Client.GetChannelByName(event.ChannelName, eh.mhHandler.Team.Id, "")
	}

	if channel != nil {
		event.ChannelID = channel.Id
		event.ChannelName = channel.Name
	}

	return event
}
