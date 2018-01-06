package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"log"
	"time"
)

type BotServer struct {
	config       Config
	slackHandler *SlackHandler
	eventHandler *EventHandler
	ruleEngine   *rules.Engine
}

func NewBotServer(config Config) *BotServer {
	b := BotServer{}
	b.config = config
	b.slackHandler = NewSlackHandler(config.BotToken, config.AdminToken)
	b.ruleEngine = rules.NewEngine(b.config.CaseDir)
	return &b
}

func (b *BotServer) Start() {
	log.Println("[+] Claptrap BotServer starting")
	rtm := b.slackHandler.StartRTM()
	attempt := 1
	for {
		botInfo := rtm.GetInfo()
		log.Printf("[+] RTM connection starting [Attempt %d]\n", attempt)
		if botInfo != nil {
			break
		}
		attempt++
		time.Sleep(time.Millisecond * 500)
	}
	if b.config.AutoJoinAllChannel {
		b.slackHandler.AutoJoinAllChannel(rtm.GetInfo().User.ID)
	}
	b.eventHandler = NewEventHandler(rtm, b.ruleEngine)
	respHandler := NewSlackResponseHandler(rtm, b.slackHandler.AdminAPI)
	b.ruleEngine.SetResponseHandler(respHandler)
	b.ruleEngine.Start()
	b.eventHandler.Start()
}

func (b *BotServer) AddCase(caseType string, newCase rules.Case) {
	b.ruleEngine.AddCase(caseType, newCase)
	log.Printf("[+] Dynamic Case '%s' with type '%s' loaded\n", newCase.Name, caseType)
}