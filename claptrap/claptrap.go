package claptrap

import (
	"github.com/DSchalla/Claptrap/rules"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path"
	"strings"
)

var caseTypes = []string{"message", "channel_join"}

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
	b.ruleEngine = rules.NewEngine()
	b.caseWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Println("[!] Unable to Create File Watcher")
	}
	return &b
}

func (b *BotServer) Start() {
	log.Println("[+] Claptrap BotServer starting")
	b.ReloadCases()
	go b.startCaseWatcher()
	rtm := b.slackHandler.StartRTM()
	b.eventHandler = NewEventHandler(rtm, b.ruleEngine)
	respHandler := NewSlackResponseHandler(rtm, b.slackHandler.AdminAPI)
	b.ruleEngine.SetResponseHandler(respHandler)
	b.eventHandler.Start()
}

func (b *BotServer) ReloadCases() {
	for _, caseType := range caseTypes {
		b.ReloadCase(caseType)
	}
}

func (b *BotServer) ReloadCase(caseType string) bool {
	valid := false
	for _, validType := range caseTypes {
		if validType == caseType {
			valid = true
			break
		}
	}

	if !valid {
		return false
	}

	filePath := path.Join(b.config.CaseDir, caseType+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	cases := loadCasesFromFile(filePath)
	b.ruleEngine.LoadCases(caseType, cases)
	log.Printf("[+] %d cases loaded for type '%s'", len(cases), caseType)
	return true
}

func (b *BotServer) startCaseWatcher() {
	err := b.caseWatcher.Add(b.config.CaseDir)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-b.caseWatcher.Events:
			if !strings.HasSuffix(event.Name, ".json") {
				continue
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				log.Println("[+] Case Config modified:", event.Name)
				caseType := strings.Replace(event.Name, ".json", "", 1)
				caseType = strings.Replace(caseType, b.config.CaseDir, "", 1)
				b.ReloadCase(caseType)
			}
		case err := <-b.caseWatcher.Errors:
			log.Println("[!] Error From Case Watcher:", err)
		}
	}

}
