package claptrap

import (
	"encoding/json"
	"fmt"
	"github.com/DSchalla/Claptrap/rules"
	"io/ioutil"
)

type rawCase struct {
	Name       string         `json:"name"`
	Conditions []rawCondition `json:"conditions"`
	Responses  []rawResponse  `json:"responses"`
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

func loadCasesFromFile(filepath string) []rules.Case {
	var rawCases []rawCase
	raw, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	json.Unmarshal(raw, &rawCases)

	cases := make([]rules.Case, len(rawCases))
	for i, rawCase := range rawCases {
		cases[i] = createCaseFromRawCase(rawCase)
	}

	return cases
}
func createCaseFromRawCase(r rawCase) rules.Case {
	parsedCase := rules.Case{}
	parsedCase.Name = r.Name
	parsedConditions := make([]rules.Condition, len(r.Conditions))
	parsedResponses := make([]rules.Response, len(r.Responses))

	for i, condition := range r.Conditions {
		parsedConditions[i] = createConditionFromRawCondition(condition)
	}

	for i, response := range r.Responses {
		parsedResponses[i] = createResponseFromRawResponse(response)
	}

	parsedCase.Conditions = parsedConditions
	parsedCase.Responses = parsedResponses

	return parsedCase
}

func createConditionFromRawCondition(rawCond rawCondition) rules.Condition {
	var realCondition rules.Condition

	switch condType := rawCond.CondType; condType {
	case "text_contains":
		realCondition = rules.TextContainsCondition{Condition: rawCond.Condition}
	case "text_equals":
		realCondition = rules.TextEqualsCondition{Condition: rawCond.Condition}
	case "user_equals":
		realCondition = rules.UserEqualsCondition{Condition: rawCond.Condition, Parameter: rawCond.Parameter}
	case "user_is_role":
		realCondition = rules.UserIsRoleCondition{Condition: rawCond.Condition, Parameter: rawCond.Parameter}
	case "channel_equals":
		realCondition = rules.ChannelEqualsCondition{Condition: rawCond.Condition}
	case "channel_is_type":
		realCondition = rules.ChannelIsTypeCondition{Condition: rawCond.Condition}
	case "random":
		realCondition = rules.RandomCondition{Likeness: rawCond.Likeness}
	default:
		fmt.Errorf("Invalid Condition Type: %s\n", condType)
	}

	return realCondition
}

func createResponseFromRawResponse(rawResp rawResponse) rules.Response {
	var realResponse rules.Response

	switch respType := rawResp.Action; respType {
	case "message_channel":
		realResponse = rules.MessageChannelResponse{ChannelID: rawResp.Channel, Message: rawResp.Message}
	case "message_user":
		realResponse = rules.MessageUserResponse{UserID: rawResp.User, Message: rawResp.Message}
	case "invite_user":
		realResponse = rules.InviteUserResponse{ChannelID: rawResp.Channel, UserID: rawResp.User}
	case "kick_user":
		realResponse = rules.KickUserResponse{ChannelID: rawResp.Channel, UserID: rawResp.User}
	case "delete_message":
		realResponse = rules.DeleteMessageResponse{}
	default:
		fmt.Errorf("Invalid Condition Type: %s\n", respType)
	}

	return realResponse
}
