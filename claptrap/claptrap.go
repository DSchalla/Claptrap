package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"log"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/model"
			"fmt"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/DSchalla/Claptrap/provider"
	"sync"
)

const RejectMessage = "Message intercepted by Rule Engine (CL4P-TP)"

type BotServer struct {
	config     Config
	api		   plugin.API
	prov 	   provider.Provider
	ruleEngine *rules.Engine
	evaluationMutex sync.RWMutex
}

func NewBotServer(api plugin.API, config Config) (*BotServer, error) {
	b := BotServer{}
	b.evaluationMutex = sync.RWMutex{}
	b.evaluationMutex.Lock()
	defer b.evaluationMutex.Unlock()
	b.config = config
	b.api = api
	botUser, err := api.GetUserByUsername(config.Name)
	if err != nil {
		return nil, err
	}
	mlog.Debug(fmt.Sprintf("BotUser configured: %s (%s)", botUser.Username, botUser.Id))
	b.prov = provider.NewMattermost(api, botUser)
	b.ruleEngine = rules.NewEngine(b.prov)

	return &b, nil
}

func (b *BotServer) HandleMessage(post *model.Post, intercept bool) (*model.Post, string) {
	b.evaluationMutex.RLock()
	defer b.evaluationMutex.RUnlock()
	mlog.Debug(fmt.Sprintf("Message received: %+v\n", post))
	event := b.prov.NormalizeMessageEvent(post)
	res := b.ruleEngine.EvaluateEvent(event, intercept)

	if res.Intercept {
		return nil, RejectMessage
	}
	return post, ""
}

func (b *BotServer) ReloadConfig(config Config) {
	b.evaluationMutex.Lock()
	defer b.evaluationMutex.Unlock()
	botUser, _ := b.api.GetUserByUsername(config.Name)
	mlog.Debug(fmt.Sprintf("BotUser configured: %s (%s)", botUser.Username, botUser.Id))
	b.prov = provider.NewMattermost(b.api, botUser)
	b.ruleEngine = rules.NewEngine(b.prov)
	cases := rules.LoadCasesFromString(b.config.Cases)

	for _, botCase := range cases {
		b.ruleEngine.AddCase("message", botCase)
	}
}

func (b *BotServer) AddCase(caseType string, newCase rules.Case) {
	b.ruleEngine.AddCase(caseType, newCase)
	log.Printf("[+] Dynamic Case '%s' with type '%s' loaded\n", newCase.Name, caseType)
}
