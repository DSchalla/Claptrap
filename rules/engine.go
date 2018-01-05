package rules

import (
	"github.com/fsnotify/fsnotify"
	"golang.org/x/tools/go/gcimporter15/testdata"
	"log"
	"os"
	"path"
	"strings"
)

var caseTypes = []string{"message", "channel_join"}

type Engine struct {
	caseDir         string
	caseFiles       map[string][]Case
	caseDynamic     map[string][]Case
	caseCombined    map[string][]Case
	responseHandler ResponseHandler
	caseWatcher     *fsnotify.Watcher
}

func NewEngine(caseDir string) *Engine {
	e := &Engine{}
	e.caseFiles = make(map[string][]Case)
	e.caseDynamic = make(map[string][]Case)
	e.caseCombined = make(map[string][]Case)
	e.caseDir = caseDir
	return e
}

func (e *Engine) Start() {
	e.ReloadCaseFiles()
	go e.startCaseFileWatcher()
}

func (e *Engine) SetResponseHandler(handler ResponseHandler) {
	e.responseHandler = handler
}

func (e *Engine) AddCase(caseType string, cases []Case) {
	e.caseFiles[caseType] = cases
	e.combineCaseMap()
}

func (e *Engine) ReloadCaseFiles() {
	for _, caseType := range caseTypes {
		e.ReloadCaseFile(caseType)
	}
}

func (e *Engine) ReloadCaseFile(caseType string) bool {
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

	filePath := path.Join(e.caseDir, caseType+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	cases := loadCasesFromFile(filePath)
	e.caseFiles[caseType] = cases
	log.Printf("[+] %d file cases loaded for type '%s'", len(cases), caseType)
	e.combineCaseMap()
	return true
}

func (e *Engine) startCaseFileWatcher() {
	err := e.caseWatcher.Add(e.caseDir)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-e.caseWatcher.Events:
			if !strings.HasSuffix(event.Name, ".json") {
				continue
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				log.Println("[+] Case Config modified:", event.Name)
				caseType := strings.Replace(event.Name, ".json", "", 1)
				caseType = strings.Replace(caseType, e.caseDir, "", 1)
				e.ReloadCaseFile(caseType)
			}
		case err := <-e.caseWatcher.Errors:
			log.Println("[!] Error From Case Watcher:", err)
		}
	}

}

func (e *Engine) combineCaseMap() {
	for _, caseType := range caseTypes {
		e.caseCombined[caseType] = append(e.caseFiles[caseType], e.caseDynamic[caseType]...)
	}
}

func (e *Engine) EvaluateEvent(event Event) bool {
	log.Printf(
		"[+] Event received of type '%s' by '%s' in '%s' \n",
		event.Type, event.UserName, event.ChannelName,
	)
	cases := e.caseCombined[event.Type]
	return e.checkCases(event, cases)
}

func (e *Engine) checkCases(event Event, cases []Case) bool {
	hitCase := false
	for _, eventCase := range cases {
		if e.checkConditions(event, eventCase.Conditions) {
			log.Printf(
				"[+] Case '%s' matched", eventCase.Name)
			e.executeResponse(event, eventCase)
			hitCase = true
		}
	}
	return hitCase
}

func (e *Engine) checkConditions(event Event, conditions []Condition) bool {
	checkResults := make([]bool, len(conditions))
	for i, condition := range conditions {
		checkResults[i] = condition.Test(event)
	}

	valid := true
	for _, checkResult := range checkResults {
		if !checkResult {
			valid = false
			break
		}
	}

	return valid
}

func (e *Engine) executeResponse(event Event, eventCase Case) bool {

	if eventCase.ResponseFunc != nil {
		return eventCase.ResponseFunc(event, e.responseHandler)
	}

	for _, response := range eventCase.Responses {
		response.Execute(e.responseHandler, event)
	}

	return true
}
