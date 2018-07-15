package rules

import (
	"github.com/DSchalla/Claptrap/provider"
	"log"
		"github.com/DSchalla/Claptrap/analysis"
	"time"
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
			e.audit.Add(analysis.CaseTriggerAuditEvent{
				Username: event.UserName,
				UserId: event.UserID,
				CaseId: eventCase.Name,
				Timestamp: time.Now(),
			})
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

		if condition == nil {
			continue
		}

		checkResults[i] = condition.Test(event)
	}

	return e.conditionMatcher.Evaluate(matching, checkResults)
}

func (e *Engine) executeResponse(event provider.Event, eventCase Case) bool {

	auditEvent := analysis.ActionExecutedAuditEvent{
		ChannelId: event.ChannelID,
		ChannelName: event.ChannelName,
		TeamId: event.TeamID,
		UserId: event.UserID,
		Username: event.UserName,
		CaseId: eventCase.Name,
		Timestamp: time.Now(),
	}

	if eventCase.ResponseFunc != nil {
		auditEvent.Action = "CustomFunc"
		e.audit.Add(auditEvent)
		return eventCase.ResponseFunc(event, e.provider)
	}

	for _, response := range eventCase.Responses {

		if response == nil {
			continue
		}
		auditEvent.Action = response.GetName()
		e.audit.Add(auditEvent)
		response.Execute(e.provider, event)
	}

	return true
}
