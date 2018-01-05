package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path"
	"strings"
)

type BotServer struct {
	config       Config
	slackHandler *SlackHandler
	eventHandler *EventHandler
	ruleEngine   *rules.Engine
	caseWatcher  *fsnotify.Watcher
}

func NewBotServer(config Config) *BotServer {
	var err error
	b := BotServer{}
	b.config = config
	b.slackHandler = NewSlackHandler(config.BotToken, config.AdminToken)
	b.ruleEngine = rules.NewEngine(b.config.CaseDir)
	b.caseWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Println("[!] Unable to Create File Watcher")
	}
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
