package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"github.com/nlopes/slack"
	"log"
	"strings"
)

type EventHandler struct {
	rtm        *slack.RTM
	ruleEngine *rules.Engine
}

func NewEventHandler(rtm *slack.RTM, ruleEngine *rules.Engine) *EventHandler {
	eh := EventHandler{
		rtm: rtm,
	}
	eh.ruleEngine = ruleEngine
	return &eh
}

func (eh *EventHandler) Start() {
	log.Println("[+] Event Handler started and listening")
	for msg := range eh.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			go eh.handleMessageEvent(ev)
		case *slack.ConnectionErrorEvent:
			log.Printf("[!] Error Connecting to Slack RTM: %s (Attempt: %d)\n", ev.ErrorObj, ev.Attempt)
		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

func (eh *EventHandler) handleMessageEvent(event *slack.MessageEvent) {

	if event.User == "" {
		//ToDo: Specify Subtypes to return early
		return
	}

	if event.User == eh.rtm.GetInfo().User.ID {
		// Ignore messages posted by the Bot itself
		return
	}

	unifiedEvent := rules.Event{}
	unifiedEvent.Type = "message"
	unifiedEvent.UserID = event.User
	unifiedEvent.ChannelID = event.Channel
	unifiedEvent.Timestamp = event.Timestamp
	unifiedEvent.Text = event.Text
	unifiedEvent = eh.addEventMetadata(unifiedEvent)
	eh.ruleEngine.EvaluateEvent(unifiedEvent)
}

func (eh *EventHandler) addEventMetadata(event rules.Event) rules.Event {
	var userName, userRole, channelName, inviterName, inviterRole string

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

	return event
}
