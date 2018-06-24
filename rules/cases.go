package rules

import (
	"encoding/json"
	"fmt"
	"github.com/DSchalla/Claptrap/provider"
	"io/ioutil"
	"log"
	"strings"
)

type Case struct {
	Name              string
	Intercept		  bool
	ConditionMatching string
	Conditions        []Condition
	Responses         []Response
	ResponseFunc      func(event provider.Event, p provider.Provider) bool
}

type rawCase struct {
	Name              string         `json:"name"`
	ConditionMatching string         `json:"condition_matching"`
	Conditions        []rawCondition `json:"conditions"`
	Responses         []rawResponse  `json:"responses"`
}

type rawCondition struct {
	CondType  string `json:"type"`
	Condition string `json:"condition"`
	Likeness  int    `json:"likeness"`
	Parameter string `json:"parameter"`
}

type rawResponse struct {
	Action  string `json:"action"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	Message string `json:"message"`
}

func loadCasesFromFile(filepath string) []Case {
	var rawCases []rawCase
	raw, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	json.Unmarshal(raw, &rawCases)

	cases := make([]Case, len(rawCases))
	for i, rawCase := range rawCases {
		cases[i] = createCaseFromRawCase(rawCase)
	}

	return cases
}

func LoadCasesFromString(caseString string) []Case {
	var rawCases []rawCase

	json.Unmarshal([]byte(caseString), &rawCases)

	cases := make([]Case, len(rawCases))
	for i, rawCase := range rawCases {
		cases[i] = createCaseFromRawCase(rawCase)
	}

	return cases
}

func createCaseFromRawCase(r rawCase) Case {
	parsedCase := Case{}
	parsedCase.Name = r.Name
	parsedCase.ConditionMatching = strings.ToLower(r.ConditionMatching)

	if parsedCase.ConditionMatching == "" {
		parsedCase.ConditionMatching = "and"
	}

	parsedConditions := make([]Condition, len(r.Conditions))
	parsedResponses := make([]Response, len(r.Responses))

	for i, condition := range r.Conditions {
		parsedCondition, err := createConditionFromRawCondition(condition)

		if err == nil {
			parsedConditions[i] = parsedCondition
		} else {
			log.Printf("[!] Error creating condition %s of case %s: %s -> Skipped\n", condition.CondType, parsedCase.Name, err)
		}
	}

	for i, response := range r.Responses {
		parsedResponse, err := createResponseFromRawResponse(response)

		if err == nil {
			parsedResponses[i] = parsedResponse
		} else {
			log.Printf("[!] Error creating response %s of case %s: %s -> Skipped\n", response.Action, parsedCase.Name, err)
		}
	}

	parsedCase.Conditions = parsedConditions
	parsedCase.Responses = parsedResponses

	return parsedCase
}

func createConditionFromRawCondition(rawCond rawCondition) (Condition, error) {
	var realCondition Condition
	var err error
	switch condType := rawCond.CondType; condType {
	case "text_contains":
		realCondition, err = NewTextContainsCondition(rawCond.Condition)
	case "text_equals":
		realCondition, err = NewTextEqualsCondition(rawCond.Condition)
	case "text_starts_with":
		realCondition, err = NewTextStartsWithCondition(rawCond.Condition)
	case "text_matches":
		realCondition, err = NewTextMatchesCondition(rawCond.Condition)
	case "user_equals":
		realCondition, err = NewUserEqualsCondition(rawCond.Condition, rawCond.Parameter)
	case "user_is_role":
		realCondition, err = NewUserIsRoleCondition(rawCond.Condition, rawCond.Parameter)
	case "channel_equals":
		realCondition, err = NewChannelEqualsCondition(rawCond.Condition)
	case "channel_is_type":
		realCondition, err = NewChannelIsTypeCondition(rawCond.Condition)
	case "random":
		realCondition, err = NewRandomCondition(rawCond.Likeness)
	default:
		err = fmt.Errorf("Invalid Condition Type: %s\n", condType)
	}

	return realCondition, err
}

func createResponseFromRawResponse(rawResp rawResponse) (Response, error) {
	var realResponse Response
	var err error
	switch respType := rawResp.Action; respType {
	case "message_channel":
		realResponse, err = NewMessageChannelResponse(rawResp.Channel, rawResp.Message)
	case "message_user":
		realResponse, err = NewMessageUserResponse(rawResp.User, rawResp.Message)
	case "invite_user":
		realResponse, err = NewInviteUserResponse(rawResp.Channel, rawResp.User)
	case "kick_user":
		realResponse, err = NewKickUserResponse(rawResp.Channel, rawResp.User)
	case "delete_message":
		realResponse, err = NewDeleteMessageResponse()
	default:
		err = fmt.Errorf("Invalid Response Type: %s\n", respType)
	}

	return realResponse, err
}
