package rules

import (
	"github.com/DSchalla/Claptrap/provider"
	"log"
		"github.com/DSchalla/Claptrap/analysis"
)

var caseTypes = []string{"Message", "user_add", "user_remove"}

type Engine struct {
	caseManager      *CaseManager
	provider         provider.Provider
	conditionMatcher *ConditionMatcher
	audit            *analysis.AuditTrail
}

type Result struct {
	Hit       bool
	HitCases  []Case
	Intercept bool
}

func NewEngine(caseManager *CaseManager, provider provider.Provider, audit *analysis.AuditTrail) *Engine {
	e := &Engine{}
	e.provider = provider
	e.caseManager = caseManager
	e.audit = audit
	e.conditionMatcher = NewConditionMatcher()
	return e
}

func (e *Engine) EvaluateEvent(event provider.Event, intercept bool) Result {
	log.Printf(
		"[+] Event received of type '%s' by '%s' in '%s' \n",
		event.Type, event.UserName, event.ChannelName,
	)
	//ToDo: Handle error gracefully and log
	cases, _ := e.caseManager.GetForType(event.Type)
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
