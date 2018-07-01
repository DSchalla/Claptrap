package analysis

import (
	"time"
	"fmt"
	"github.com/mattermost/mattermost-server/plugin"
	"bytes"
	"encoding/gob"
	"sync"
)

func NewAuditTrail(api plugin.API) *AuditTrail {
	a := &AuditTrail{}
	a.api = api
	a.mutex = &sync.RWMutex{}

	gob.Register(CaseTriggerAuditEvent{})
	gob.Register(CaseCreatedAuditEvent{})
	gob.Register(CaseDeletedAuditEvent{})
	gob.Register(ActionExecutedAuditEvent{})

	return a
}

type AuditTrail struct {
	api   plugin.API
	mutex *sync.RWMutex
}

func (a *AuditTrail) Add(message AuditMessage) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	messages, err := a.GetEvents(message.GetTimestamp())

	if err != nil {
		return err
	}

	messages = append(messages, message)

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(messages)

	if err != nil {
		return err
	}

	data := buffer.Bytes()
	key := a.getKey(message.GetTimestamp())
	a.api.KVSet(key, data)

	return nil
}

func (a *AuditTrail) GetEvents(timestamp time.Time) ([]AuditMessage, error) {
	var buffer bytes.Buffer
	var messages []AuditMessage

	a.mutex.RLock()
	defer a.mutex.RUnlock()

	dec := gob.NewDecoder(&buffer)
	key := a.getKey(timestamp)
	data, err := a.api.KVGet(key)

	if err != nil {
		return nil, err
	}

	if data != nil {
		buffer.Write(data)
		err2 := dec.Decode(&messages)
		if err2 != nil {
			return nil, err2
		}
	}

	return messages, nil
}

func (a *AuditTrail) getKey(t time.Time) string {
	key := t.Format("2006-01-02")
	return "audit." + key
}

type AuditMessage interface {
	String() string
	GetTimestamp() time.Time
}

type CaseTriggerAuditEvent struct {
	Username  string
	UserId    string
	CaseId    string
	Timestamp time.Time
}

func (c CaseTriggerAuditEvent) String() string {
	return fmt.Sprintf("User '%s' (%s) triggered case '%s' at %s", c.Username, c.UserId, c.CaseId, c.Timestamp)
}

func (c CaseTriggerAuditEvent) GetTimestamp() time.Time {
	return c.Timestamp
}

type CaseCreatedAuditEvent struct {
	Username  string
	UserId    string
	CaseId    string
	Timestamp time.Time
}

func (c CaseCreatedAuditEvent) String() string {
	return fmt.Sprintf("User '%s' (%s) created case '%s' at %s", c.Username, c.UserId, c.CaseId, c.Timestamp)
}

func (c CaseCreatedAuditEvent) GetTimestamp() time.Time {
	return c.Timestamp
}

type CaseDeletedAuditEvent struct {
	Username  string
	UserId    string
	CaseId    string
	Timestamp time.Time
}

func (c CaseDeletedAuditEvent) String() string {
	return fmt.Sprintf("User '%s' (%s) deleted case '%s' at %s", c.Username, c.UserId, c.CaseId, c.Timestamp)
}

func (c CaseDeletedAuditEvent) GetTimestamp() time.Time {
	return c.Timestamp
}

type ActionExecutedAuditEvent struct {
	Action      string
	Username    string
	UserId      string
	ChannelName string
	ChannelId   string
	CaseId      string
	Timestamp   time.Time
}

func (c ActionExecutedAuditEvent) String() string {

	log := "Executed action '%s' "

	if c.ChannelId == "" {
		log += fmt.Sprintf("against Channel '%s' (%s)", c.ChannelName, c.ChannelId)
	} else {
		log += fmt.Sprintf("against User '%s' (%s)", c.Username, c.UserId)
	}

	log += fmt.Sprintf(" as part of case '%s' at %s", c.CaseId, c.Timestamp)

	return log
}

func (c ActionExecutedAuditEvent) GetTimestamp() time.Time {
	return c.Timestamp
}
