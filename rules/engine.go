package rules

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path"
	"strings"
	"path/filepath"
)

var caseTypes = []string{"message", "user_add", "user_remove"}

type Engine struct {
	caseDir         string
	caseFiles       map[string][]Case
	caseDynamic     map[string][]Case
	caseCombined    map[string][]Case
	responseHandler ResponseHandler
	caseWatcher     *fsnotify.Watcher
}

func NewEngine(caseDir string) *Engine {
	var err error
	e := &Engine{}
	e.caseFiles = make(map[string][]Case)
	e.caseDynamic = make(map[string][]Case)
	e.caseCombined = make(map[string][]Case)
	e.caseDir = caseDir
	e.caseWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Println("[!] Unable to Create File Watcher")
	}
	return e
}

func (e *Engine) Start() {
	e.ReloadCaseFiles()
	go e.startCaseFileWatcher()
}

func (e *Engine) SetResponseHandler(handler ResponseHandler) {
	e.responseHandler = handler
}

func (e *Engine) AddCase(caseType string, newCase Case) {
	e.caseDynamic[caseType] = append(e.caseDynamic[caseType], newCase)
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
		log.Printf("[+] Invalid Case Type '%s'\n", caseType)
		return false
	}

	filePath := path.Join(e.caseDir, caseType+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("[+] File '%s' does not exist\n", filePath)
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
				caseType := filepath.Base(event.Name)
				caseType = strings.Replace(caseType, ".json", "", 1)
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
		} else {
			log.Printf("[+] Case '%s' did not match", eventCase.Name)
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
