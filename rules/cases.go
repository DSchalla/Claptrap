package rules

import (
	"fmt"
	"github.com/DSchalla/Claptrap/provider"
	"log"
	"strings"
	"sync"
	"errors"
	"time"
	"github.com/mattermost/mattermost-server/plugin"
	"bytes"
	"encoding/gob"
	"github.com/mattermost/mattermost-server/mlog"
)

type Case struct {
	Name              string
	Intercept         bool
	ConditionMatching string
	Conditions        []Condition
	Responses         []Response
	ResponseFunc      func(event provider.Event, p provider.Provider) bool
	DeleteTime        time.Time
}

var ValidTypes = []string{"message", "channel_join"}

func NewCaseManager(api plugin.API) *CaseManager {
	cm := &CaseManager{}
	cm.mutex = &sync.RWMutex{}
	cm.api = api

	gob.Register(TextContainsCondition{})
	gob.Register(TextEqualsCondition{})
	gob.Register(TextStartsWithCondition{})
	gob.Register(TextMatchesCondition{})
	gob.Register(UserIsRoleCondition{})
	gob.Register(MessageUserResponse{})
	gob.Register(MessageChannelResponse{})
	gob.Register(DeleteMessageResponse{})

	return cm
}

type CaseManager struct {
	mutex *sync.RWMutex
	api   plugin.API
}

func (c *CaseManager) Add(caseType string, newCase Case) error {
	var buffer bytes.Buffer
	var cases []Case

	if !c.validType(caseType) {
		return errors.New("invalid case type")
	}

	cases, err := c.GetForType(caseType)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err != nil {
		return err
	}

	cases = append(cases, newCase)
	fmt.Printf("%+v", cases)
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(cases)

	if err != nil {
		return err
	}

	data := buffer.Bytes()
	c.api.KVSet("cases."+caseType, data)

	fmt.Printf("%+v", data)

	return nil
}

func (c *CaseManager) Delete(caseType, caseName string) error {
	var buffer bytes.Buffer
	var cases []Case

	if !c.validType(caseType) {
		return errors.New("invalid case type")
	}

	cases, err := c.GetForType(caseType)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err != nil {
		return err
	}

	var newCases []Case

	for _, currentCase := range cases {
		if currentCase.Name == caseName {
			continue
		}

		newCases = append(newCases, currentCase)
	}

	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(newCases)

	if err != nil {
		return err
	}

	data := buffer.Bytes()
	c.api.KVSet("cases."+caseType, data)

	fmt.Printf("%+v", data)

	return nil
}

func (c *CaseManager) GetForType(caseType string) ([]Case, error) {
	var buffer bytes.Buffer
	var cases []Case

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	dec := gob.NewDecoder(&buffer)
	data, err := c.api.KVGet("cases." + caseType)

	if err != nil {
		mlog.Debug(fmt.Sprintf("[CLAPTRAP-PLUGIN][CaseManager] Unable to get cases: %s", err))
		return nil, err
	}

	if data != nil {
		buffer.Write(data)
		err2 := dec.Decode(&cases)
		if err2 != nil {
			mlog.Debug(fmt.Sprintf("[CLAPTRAP-PLUGIN][CaseManager] Unable to decode cases: %s", err2))
			return nil, err2
		}
	}

	fmt.Printf("%+v", cases)

	return cases, nil
}

func (c *CaseManager) GetCaseTypes() map[string]string {
	return map[string]string{
		"message": "Message",
		"channel_join": "Channel Join (incl. Invite)",
		"channel_leave": "Channel Leave (incl. Kick)",
		"team_join": "Team Join",
	}
}

func (c *CaseManager) GetConditionOptions() map[string]string{
	return map[string]string{
	 "message_contains": "Message Contains",
	 "message_equals": "Message Equals",
	 "message_starts_with": "Message Starts With",
	 "message_matches": "Message Matches",
	 "user_equals": "User Equals",
	 "user_is_role": "User Is Role",
	 "channel_equals": "Channel Equals",
	 "channel_is_type": "Channel Is Type",
	 "random": "Random",
	}
}

func (c *CaseManager) GetResponseOptions() map[string]string{
	return map[string]string{
	 "message_channel": "Message Channel",
	 "message_user": "Message User",
	 "message_ephemeral": "Message Ephemeral",
	 "invite_user": "Invite User",
	 "kick_user": "Kick User",
	 "delete_message": "Delete Message",
	 "intercept": "Intercept Event",
	}
}

func (c *CaseManager) validType(caseType string) bool {
	valid := false

	for _, validType := range ValidTypes {
		if validType == caseType {
			valid = true
			break
		}
	}

	return valid
}

type RawCase struct {
	Name              string         `json:"name"`
	ConditionMatching string         `json:"condition_matching"`
	Intercept bool         `json:"intercept"`
	Conditions        []RawCondition `json:"conditions"`
	Responses         []RawResponse  `json:"responses"`
}

type RawCondition struct {
	CondType  string `json:"type"`
	Condition string `json:"condition"`
	Likeness  int    `json:"likeness"`
	Parameter string `json:"parameter"`
}

type RawResponse struct {
	Action  string `json:"action"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	Message string `json:"Message"`
}

func CreateCaseFromRawCase(r RawCase) Case {
	parsedCase := Case{}
	parsedCase.Name = r.Name
	parsedCase.ConditionMatching = strings.ToLower(r.ConditionMatching)
	parsedCase.Intercept = r.Intercept

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

func createConditionFromRawCondition(rawCond RawCondition) (Condition, error) {
	var realCondition Condition
	var err error
	switch condType := rawCond.CondType; condType {
	case "message_contains":
		realCondition, err = NewTextContainsCondition(rawCond.Condition)
	case "message_equals":
		realCondition, err = NewTextEqualsCondition(rawCond.Condition)
	case "message_starts_with":
		realCondition, err = NewTextStartsWithCondition(rawCond.Condition)
	case "message_matches":
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

func createResponseFromRawResponse(rawResp RawResponse) (Response, error) {
	var realResponse Response
	var err error
	switch respType := rawResp.Action; respType {
	case "message_channel":
		realResponse, err = NewMessageChannelResponse(rawResp.Channel, rawResp.Message)
	case "message_user":
		realResponse, err = NewMessageUserResponse(rawResp.User, rawResp.Message)
	case "message_ephemeral":
		realResponse, err = NewMessageEphemeralResponse(rawResp.Message)
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
