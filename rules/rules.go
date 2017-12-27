package rules

import (
	"log"
)

type Engine struct {
	cases           map[string][]Case
	responseHandler ResponseHandler
}

func NewEngine() *Engine {
	e := &Engine{}
	e.cases = make(map[string][]Case)
	return e
}

func (e *Engine) SetResponseHandler(handler ResponseHandler) {
	e.responseHandler = handler
}

func (e *Engine) LoadCases(caseType string, cases []Case) {
	e.cases[caseType] = cases
}

func (e *Engine) EvaluateEvent(event Event) bool {
	log.Printf(
		"[+] Event received of type '%s' by '%s' in '%s' \n",
		event.Type, event.UserName, event.ChannelName,
	)
	cases := e.cases[event.Type]
	return e.checkCases(event, cases)
}

func (e *Engine) checkCases(event Event, cases []Case) bool {
	hitCase := false
	for _, eventCase := range cases {
		if e.checkConditions(event, eventCase.Conditions) {
			log.Printf(
				"[+] Case '%s' matched", eventCase.Name)
			e.executeResponses(event, eventCase.Responses)
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

func (e *Engine) executeResponses(event Event, responses []Response) bool {
	for _, response := range responses {
		response.Execute(e.responseHandler, event)
	}
	return true
}
