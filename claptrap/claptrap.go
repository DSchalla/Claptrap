package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"log"
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