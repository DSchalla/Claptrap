package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/model"
	"fmt"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/DSchalla/Claptrap/provider"
	"sync"
	"net/http"
	"github.com/DSchalla/Claptrap/web"
	"github.com/DSchalla/Claptrap/analysis"
)

const RejectMessage = "Message intercepted by Rule Engine (CL4P-TP)"

type Config struct {
	Name               string
	AutoJoinAllChannel bool
}

type BotServer struct {
	config Config
	audit  *analysis.AuditTrail

	api  plugin.API
	prov provider.Provider

	caseManager     *rules.CaseManager
	ruleEngine      *rules.Engine
	evaluationMutex sync.RWMutex

	webServer *web.Server
}

func NewBotServer(api plugin.API, config Config) (*BotServer, error) {
	mlog.Debug("Bot Server Started")
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

	b.audit = analysis.NewAuditTrail(api)
	b.audit.LogStart()
	b.caseManager = rules.NewCaseManager(api)
	b.prov = provider.NewMattermost(api, botUser)
	b.ruleEngine = rules.NewEngine(b.caseManager, b.prov, b.audit)
	b.webServer = web.NewServer(api, b.caseManager, b.audit)

	return &b, nil
}

func (b *BotServer) Shutdown() error {
	b.audit.LogShutdown()
	return nil
}

func (b *BotServer) HandleMessage(post *model.Post, intercept bool) (*model.Post, string) {
	mlog.Debug(fmt.Sprintf("Message received: %+v\n", post))

	b.evaluationMutex.RLock()
	defer b.evaluationMutex.RUnlock()

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
	b.ruleEngine = rules.NewEngine(b.caseManager, b.prov, b.audit)
}

func (b *BotServer) AddCase(caseType string, newCase rules.Case) {
	b.caseManager.Add(caseType, newCase)
	mlog.Debug(fmt.Sprintf("[+] Dynamic Case '%s' with type '%s' loaded\n", newCase.Name, caseType))
}

func (b *BotServer) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	b.webServer.HandleHTTP(w, r)
}
