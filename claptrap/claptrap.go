package claptrap

import (
	"github.com/DSchalla/Claptrap/provider"
	"github.com/DSchalla/Claptrap/rules"
	"log"
)

type BotServer struct {
	config     Config
	provider   provider.Provider
	ruleEngine *rules.Engine
}

func NewBotServer(config Config) (*BotServer, error) {
	b := BotServer{}
	b.config = config
	b.provider = provider.NewMattermost(config.Mattermost.ApiUrl, config.Mattermost.Username, config.Mattermost.Password, config.Mattermost.Team)
	b.ruleEngine = rules.NewEngine(b.config.General.CaseDir)
	return &b, nil
}

func (b *BotServer) Start() {
	log.Println("[+] Claptrap BotServer starting")
	b.provider.Connect()

	if b.config.General.AutoJoinAllChannel {
		go b.provider.AutoJoinAllChannel()
	}
	b.ruleEngine.SetProvider(b.provider)
	b.ruleEngine.Start()
	go b.provider.ListenForEvents()

	for event := range b.provider.GetEvents() {
		go b.ruleEngine.EvaluateEvent(event)
	}
}

func (b *BotServer) AddCase(caseType string, newCase rules.Case) {
	b.ruleEngine.AddCase(caseType, newCase)
	log.Printf("[+] Dynamic Case '%s' with type '%s' loaded\n", newCase.Name, caseType)
}
