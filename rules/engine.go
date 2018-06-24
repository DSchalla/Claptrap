package rules

import (
	"github.com/DSchalla/Claptrap/provider"
	"log"
	"os"
	"path"
)

var caseTypes = []string{"message", "user_add", "user_remove"}

type Engine struct {
	caseDir      string
	caseFiles    map[string][]Case
	caseDynamic  map[string][]Case
	caseCombined map[string][]Case
	provider	provider.Provider
	conditionMatcher *ConditionMatcher
}

type Result struct {
	Hit bool
	HitCases []Case
	Intercept bool
}

func NewEngine(provider provider.Provider) *Engine {
	e := &Engine{}
	e.provider = provider
	e.caseFiles = make(map[string][]Case)
	e.caseDynamic = make(map[string][]Case)
	e.caseCombined = make(map[string][]Case)
	e.conditionMatcher = NewConditionMatcher()
	return e
}

func (e *Engine) Start() {
	e.ReloadCaseFiles()
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

func (e *Engine) combineCaseMap() {
	for _, caseType := range caseTypes {
		e.caseCombined[caseType] = append(e.caseFiles[caseType], e.caseDynamic[caseType]...)
	}
}

func (e *Engine) EvaluateEvent(event provider.Event, intercept bool) Result {
	log.Printf(
		"[+] Event received of type '%s' by '%s' in '%s' \n",
		event.Type, event.UserName, event.ChannelName,
	)
	cases := e.caseCombined[event.Type]
	return e.checkCases(event, cases, intercept)
}

func (e *Engine) checkCases(event provider.Event, cases []Case, intercept bool) Result {
	res := Result{}

	for _, eventCase := range cases {

		if eventCase.Intercept != intercept {
			continue
		}

		if e.checkConditions(event, eventCase.ConditionMatching, eventCase.Conditions) {
			log.Printf(
				"[+] Case '%s' matched", eventCase.Name)
			e.executeResponse(event, eventCase)

			res.Hit = true
			res.HitCases = append(res.HitCases, eventCase)
			if eventCase.Intercept {
				res.Intercept = true
			}
		}
	}

	return res
}

func (e *Engine) checkConditions(event provider.Event, matching string, conditions []Condition) bool {

	if len(conditions) == 0 {
		return true
	}

	checkResults := make([]bool, len(conditions))
	for i, condition := range conditions {
		checkResults[i] = condition.Test(event)
	}

	return e.conditionMatcher.Evaluate(matching, checkResults)
}


func (e *Engine) executeResponse(event provider.Event, eventCase Case) bool {

	if eventCase.ResponseFunc != nil {
		return eventCase.ResponseFunc(event, e.provider)
	}

	for _, response := range eventCase.Responses {
		response.Execute(e.provider, event)
	}

	return true
}
