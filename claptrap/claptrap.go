package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"log"
	"os"
	"path"
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
	b.ruleEngine = rules.NewEngine()
	return &b
}

func (b *BotServer) Start() {
	log.Println("[+] Claptrap BotServer starting")
	b.ReloadCases()
	rtm := b.slackHandler.StartRTM()
	b.eventHandler = NewEventHandler(rtm, b.ruleEngine)
	respHandler := NewSlackResponseHandler(rtm, b.slackHandler.AdminAPI)
	b.ruleEngine.SetResponseHandler(respHandler)
	b.eventHandler.Start()
}

func (b *BotServer) ReloadCases() {
	caseTypes := []string{"message", "channel_join"}

	for _, caseType := range caseTypes {
		filePath := path.Join(b.config.ConfigDir, "cases", caseType+".json")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}
		cases := loadCasesFromFile(filePath)
		b.ruleEngine.LoadCases(caseType, cases)
		log.Printf("[+] %d cases loaded for type '%s'", len(cases), caseType)
	}
}
