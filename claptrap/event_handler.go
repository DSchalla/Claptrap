package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"log"
	"github.com/mattermost/mattermost-server/model"
	"fmt"
	"strings"
)

type EventHandler struct {
	eventChan        <- chan *model.WebSocketEvent
	ruleEngine *rules.Engine
}

func NewEventHandler(eventChan <- chan *model.WebSocketEvent, ruleEngine *rules.Engine) *EventHandler {
	eh := EventHandler{
		eventChan: eventChan,
	}
	eh.ruleEngine = ruleEngine
	return &eh
}

func (eh *EventHandler) Start() {
	log.Println("[+] Event Handler started and listening")
	for msg := range eh.eventChan {
		switch msg.Event {
			case "posted": eh.handleMessageEvent(msg)

		}
		fmt.Printf("Unexpected: %+v\n", msg)
	}
}

func (eh *EventHandler) handleMessageEvent(event *model.WebSocketEvent) {

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	unifiedEvent := rules.Event{}
	unifiedEvent.Type = "message"
	unifiedEvent.PostID = post.Id
	unifiedEvent.UserID = post.UserId
	unifiedEvent.UserName = event.Data["sender_name"].(string)
	unifiedEvent.ChannelID = post.ChannelId
	unifiedEvent.ChannelName = event.Data["channel_name"].(string)
	unifiedEvent.Timestamp = post.CreateAt
	unifiedEvent.Text = post.Message
	unifiedEvent = eh.addEventMetadata(unifiedEvent)
	eh.ruleEngine.EvaluateEvent(unifiedEvent)
}

func (eh *EventHandler) addEventMetadata(event rules.Event) rules.Event {
	/*	var userName, userRole, channelName, inviterName, inviterRole string

		user, err := eh.rtm.GetUserInfo(event.UserID)
		if err != nil {
			log.Printf("[!] error occured fetching username: %s\n", err)
			userName = ""
			userRole = ""
		} else {
			userName = user.Name
			if user.IsAdmin || user.IsPrimaryOwner || user.IsOwner {
				userRole = "admin"
			} else {
				userRole = "user"
			}
		}
		event.UserName = userName
		event.UserRole = userRole

		if strings.HasPrefix(event.ChannelID, "C") {
			channel, err := eh.rtm.GetChannelInfo(event.ChannelID)
			if err != nil {
				log.Printf("[!] error occured fetching channel info: %s\n", err)
				channelName = ""
			} else {
				channelName = channel.Name
			}
		} else if strings.HasPrefix(event.ChannelID, "G") {
			group, err := eh.rtm.GetGroupInfo(event.ChannelID)
			if err != nil {
				log.Printf("[!] error occured fetching group info: %s\n", err)
				channelName = ""
			} else {
				channelName = group.Name
			}
		} else {
			channelName = "Private Message"
		}

		event.ChannelName = channelName

		if event.InviterID != "" {
			inviter, err := eh.rtm.GetUserInfo(event.InviterID)
			if err != nil {
				log.Printf("[!] error occured fetching username: %s\n", err)
				inviterName = ""
				inviterRole = ""
			} else {
				inviterName = inviter.Name
				if inviter.IsAdmin || inviter.IsPrimaryOwner || inviter.IsOwner {
					inviterRole = "admin"
				} else {
					inviterRole = "user"
				}
			}
			event.InviterName = inviterName
			event.InviterRole = inviterRole
		}
	*/
	return event
}
