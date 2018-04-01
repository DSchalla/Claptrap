package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"log"
)

type BotServer struct {
	config       Config
	mattermostHandler *MattermostHandler
	eventHandler *EventHandler
	ruleEngine   *rules.Engine
}

func NewBotServer(config Config) *BotServer {
	b := BotServer{}
	b.config = config
	b.mattermostHandler = NewMattermostHandler(config.Mattermost.ApiUrl, config.Mattermost.Username, config.Mattermost.Password, config.Mattermost.Team)
	b.ruleEngine = rules.NewEngine(b.config.General.CaseDir)
	return &b
}

func (b *BotServer) Start() {
	log.Println("[+] Claptrap BotServer starting")
	b.mattermostHandler.StartWS()

	if b.config.General.AutoJoinAllChannel {
		go b.mattermostHandler.AutoJoinAllChannel()
	}
	b.eventHandler = NewEventHandler(b.mattermostHandler, b.ruleEngine)
	respHandler := NewMattermostResponseHandler(b.mattermostHandler.Client, b.mattermostHandler.BotUser)
	b.ruleEngine.SetResponseHandler(respHandler)
	b.ruleEngine.Start()
	b.eventHandler.Start()
}

func (b *BotServer) AddCase(caseType string, newCase rules.Case) {
	b.ruleEngine.AddCase(caseType, newCase)
	log.Printf("[+] Dynamic Case '%s' with type '%s' loaded\n", newCase.Name, caseType)
}