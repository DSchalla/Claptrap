package rules

import (
	"fmt"
	"github.com/DSchalla/Claptrap/provider"
	"strings"
	"sync"
	"errors"
	"time"
	"github.com/mattermost/mattermost-server/plugin"
	"bytes"
	"encoding/gob"
	"github.com/mattermost/mattermost-server/mlog"
	"net/http"
	"strconv"
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

type RawCase struct {
	Name              string
	ConditionMatching string
	Conditions        []RawCondition
	Responses         []RawResponse
}

type RawCondition struct {
	Id        int
	CondType  string
	Condition string
	Likeness  int
	Parameter string
}

type RawResponse struct {
	Id      int
	Action  string
	User    string
	Channel string
	Message string
}

var ValidTypes = []string{"message", "channel_join"}

func NewCaseManager(api plugin.API) *CaseManager {
	cm := &CaseManager{}
	cm.mutex = &sync.RWMutex{}
	cm.api = api

	gob.Register(ChannelEqualsCondition{})
	gob.Register(ChannelIsTypeCondition{})
	gob.Register(TextContainsCondition{})
	gob.Register(TextEqualsCondition{})
	gob.Register(TextStartsWithCondition{})
	gob.Register(TextMatchesCondition{})
	gob.Register(UserEqualsCondition{})
	gob.Register(UserIsRoleCondition{})
	gob.Register(RandomCondition{})

	gob.Register(MessageUserResponse{})
	gob.Register(MessageChannelResponse{})
	gob.Register(MessageEphemeralResponse{})
	gob.Register(KickUserResponse{})
	gob.Register(InviteUserResponse{})
	gob.Register(DeleteMessageResponse{})
	gob.Register(InterceptEventResponse{})

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

func (c *CaseManager) Exists(caseType, caseName string) bool {
	return false
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

func (c *CaseManager) GetCase(caseType, caseName string) (Case, error) {
	cases, err := c.GetForType(caseType)

	if err != nil {
		return Case{}, err
	}

	for _, storedCase := range cases {
		if storedCase.Name == caseName {
			return storedCase, nil
		}
	}

	return Case{}, errors.New("case not found")
}

func (c *CaseManager) GetCaseTypes() map[string]string {
	return map[string]string{
		"message":       "Message",
		"channel_join":  "Channel Join (incl. Invite)",
		"channel_leave": "Channel Leave (incl. Kick)",
		"team_join":     "Team Join",
	}
}

func (c *CaseManager) GetConditionOptions() map[string]string {
	return map[string]string{
		"message_contains":    "Message Contains",
		"message_equals":      "Message Equals",
		"message_starts_with": "Message Starts With",
		"message_matches":     "Message Matches",
		"user_equals":         "User Equals",
		"user_is_role":        "User Is Role",
		"channel_equals":      "Channel Equals",
		"channel_is_type":     "Channel Is Type",
		"random":              "Random",
	}
}

func (c *CaseManager) GetResponseOptions() map[string]string {
	return map[string]string{
		"message_channel":   "Message Channel",
		"message_user":      "Message User",
		"message_ephemeral": "Message Ephemeral",
		"invite_user":       "Invite User",
		"kick_user":         "Kick User",
		"delete_message":    "Delete Message",
		"intercept":         "Intercept Event",
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

func (c *CaseManager) CreateCaseFromHTTPReq(req *http.Request) (Case, string, error) {
	caseType := req.FormValue("type")
	validType := false

	for _, existingType := range c.GetCaseTypes() {
		if caseType == existingType {
			validType = true
			break
		}
	}

	if validType {
		return Case{}, "", errors.New("invalid Case Type")
	}

	if req.FormValue("casename") == "" && c.Exists(caseType, req.FormValue("casename")) {
		return Case{}, "", errors.New("case with that name already exists or empty name passed")
	}

	rawCase := RawCase{
		Name:              req.FormValue("casename"),
		ConditionMatching: req.FormValue("condition_matching"),
	}

	for i := 0; i < 10; i++ {
		prefix := fmt.Sprintf("conditions[%d]", i)
		conditionType := req.FormValue(prefix + "[type]")

		if conditionType == "" {
			break
		}

		conditionValue := req.FormValue(prefix + "[condition]")
		parameterValue := req.FormValue(prefix + "[parameter]")
		likenessValue, err := strconv.Atoi(req.FormValue(prefix + "[likeness]"))

		if err != nil {
			likenessValue = 0
		}

		rawCond := RawCondition{
			Id:        i,
			CondType:  conditionType,
			Condition: conditionValue,
			Likeness:  likenessValue,
			Parameter: parameterValue,
		}
		rawCase.Conditions = append(rawCase.Conditions, rawCond)
	}

	for i := 0; i < 10; i++ {
		prefix := fmt.Sprintf("responses[%d]", i)
		responseType := req.FormValue(prefix + "[type]")

		if responseType == "" {
			break
		}

		messageValue := req.FormValue(prefix + "[message]")
		userValue := req.FormValue(prefix + "[user]")
		channelValue := req.FormValue(prefix + "[channel]")

		rawResp := RawResponse{
			Id:      i,
			Action:  responseType,
			Message: messageValue,
			User:    userValue,
			Channel: channelValue,
		}
		rawCase.Responses = append(rawCase.Responses, rawResp)
	}

	realCase := c.CreateCaseFromRawCase(rawCase)

	return realCase, caseType, nil
}

func (c *CaseManager) CreateCaseFromRawCase(r RawCase) Case {
	parsedCase := Case{}
	parsedCase.Name = r.Name
	parsedCase.ConditionMatching = strings.ToLower(r.ConditionMatching)

	if parsedCase.ConditionMatching == "" {
		parsedCase.ConditionMatching = "and"
	}

	parsedConditions := make([]Condition, len(r.Conditions))
	parsedResponses := make([]Response, len(r.Responses))

	for i, condition := range r.Conditions {
		parsedCondition, err := c.createConditionFromRawCondition(condition)

		if err == nil {
			parsedConditions[i] = parsedCondition
		} else {
			mlog.Warn("[!] Error creating condition %s of case %s: %s -> Skipped\n",
				mlog.String("ConditionType", condition.CondType),
				mlog.String("CaseName", parsedCase.Name),
				mlog.Err(err),
			)
		}
	}

	for i, response := range r.Responses {
		parsedResponse, err := c.createResponseFromRawResponse(response)

		if err == nil {
			parsedResponses[i] = parsedResponse
		} else {
			mlog.Warn("[!] Error creating response %s of case %s: %s -> Skipped\n",
				mlog.String("Action", response.Action),
				mlog.String("CaseName", parsedCase.Name),
				mlog.Err(err),
			)
		}
	}

	parsedCase.Conditions = parsedConditions
	parsedCase.Responses = parsedResponses

	return parsedCase
}

func (c *CaseManager) CreateRawCaseFromCase(realCase Case) (*RawCase, error) {
	rawCase := &RawCase{
		Name: realCase.Name,
		ConditionMatching: realCase.ConditionMatching,
	}

	parsedConditions := make([]RawCondition, len(realCase.Conditions))
	parsedResponses := make([]RawResponse, len(realCase.Responses))

	for i, cond := range realCase.Conditions {
		rawCond := RawCondition{}
		rawCond.Id = i

		switch t := cond.(type) {
		case ChannelIsTypeCondition:
			rawCond.CondType = "channel_is_type"
			rawCond.Condition = t.Condition
		case ChannelEqualsCondition:
			rawCond.CondType = "channel_equals"
			rawCond.Condition = t.Condition
		case TextEqualsCondition:
			rawCond.CondType = "message_equals"
			rawCond.Condition = t.Condition
		case TextMatchesCondition:
			rawCond.CondType = "message_matches"
			rawCond.Condition = t.expression
		case TextStartsWithCondition:
			rawCond.CondType = "message_starts_with"
			rawCond.Condition = t.Condition
		case TextContainsCondition:
			rawCond.CondType = "message_contains"
			rawCond.Condition = t.Condition
		case UserIsRoleCondition:
			rawCond.CondType = "user_is_role"
			rawCond.Condition = t.Condition
			rawCond.Parameter = t.Parameter
		case UserEqualsCondition:
			rawCond.CondType = "user_equals"
			rawCond.Condition = t.Condition
			rawCond.Parameter = t.Parameter
		case RandomCondition:
			rawCond.CondType = "channel_equals"
			rawCond.Likeness = t.Likeness
		}

		parsedConditions = append(parsedConditions, rawCond)
	}

	for i, resp := range realCase.Responses {
		rawResp := RawResponse{}
		rawResp.Id = i

		switch t := resp.(type) {
		case MessageChannelResponse:
			rawResp.Action = "message_channel"
			rawResp.Message = t.Message
			rawResp.Channel = t.ChannelID
		case MessageUserResponse:
			rawResp.Action = "message_user"
			rawResp.Message = t.Message
			rawResp.User = t.UserID
		case MessageEphemeralResponse:
			rawResp.Action = "message_ephemeral"
			rawResp.Message = t.Message
		case InviteUserResponse:
			rawResp.Action = "invite_user"
			rawResp.User = t.UserID
			rawResp.Channel = t.ChannelID
		case KickUserResponse:
			rawResp.Action = "kick_user"
			rawResp.User = t.UserID
			rawResp.Channel = t.ChannelID
		case InterceptEventResponse:
			rawResp.Action = "intercept_event"
		case DeleteMessageResponse:
			rawResp.Action = "delete_message"
		}

		parsedResponses = append(parsedResponses, rawResp)
	}

	rawCase.Conditions = parsedConditions
	rawCase.Responses = parsedResponses

	return rawCase, nil
}

func (c *CaseManager) createConditionFromRawCondition(rawCond RawCondition) (Condition, error) {
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

func (c *CaseManager) createResponseFromRawResponse(rawResp RawResponse) (Response, error) {
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
	case "intercept_event":
		realResponse, err = NewInterceptEventResponse()
	default:
		err = fmt.Errorf("Invalid Response Type: %s\n", respType)
	}

	return realResponse, err
}
