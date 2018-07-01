package analysis

import (
	"testing"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"encoding/gob"
	"bytes"
	"time"
	"reflect"
)

func TestAuditTrail_GetEvents(t *testing.T) {
	var buffer bytes.Buffer
	timestamp := time.Date(2018, time.May, 05, 00, 00, 00, 0, time.UTC)

	messages := []AuditMessage{
		CaseTriggerAuditEvent{
			Username:  "BotUser",
			UserId:    "bot1",
			CaseId:    "TestCase",
			Timestamp: timestamp,
		},
		CaseCreatedAuditEvent{
			Username:  "User1",
			UserId:    "user1",
			CaseId:    "TestCase2",
			Timestamp: timestamp,
		},
	}
	gob.Register(CaseTriggerAuditEvent{})
	gob.Register(CaseCreatedAuditEvent{})
	enc := gob.NewEncoder(&buffer)
	enc.Encode(messages)

	api := &plugintest.API{}
	api.On("KVGet", "audit.2018-05-05").Return(buffer.Bytes(), nil)

	audit := NewAuditTrail(api)
	givenMessages, err := audit.GetEvents(timestamp)

	if err != nil {
		t.Fatalf("GetEvents returned error: %s", err)
		return
	}

	if len(givenMessages) == 0 {
		t.Fatalf("Given Messages got length of 0")
		return
	}

	givenMessage := givenMessages[0].(CaseTriggerAuditEvent)
	if !reflect.DeepEqual(givenMessage, messages[0]) {
		t.Fatalf("Given %+v, expected %+v", givenMessage, messages[0])
	}

	givenMessage2 := givenMessages[1].(CaseCreatedAuditEvent)
	if !reflect.DeepEqual(givenMessage2, messages[1]) {
		t.Fatalf("Given %+v, expected %+v", givenMessage2, messages[1])
	}
}

func TestAuditTrail_Add(t *testing.T) {
	var buffer bytes.Buffer

	timestamp := time.Date(2018, time.May, 05, 00, 00, 00, 0, time.UTC)

	event := CaseTriggerAuditEvent{
		Username:  "BotUser",
		UserId:    "bot1",
		CaseId:    "TestCase",
		Timestamp: timestamp,
	}

	messages := []AuditMessage{
		event,
	}
	gob.Register(CaseTriggerAuditEvent{})
	enc := gob.NewEncoder(&buffer)
	enc.Encode(messages)

	api := &plugintest.API{}
	api.On("KVSet", "audit.2018-05-05", buffer.Bytes()).Return(nil)
	api.On("KVGet", "audit.2018-05-05").Return(buffer.Bytes(), nil)

	audit := NewAuditTrail(api)
	err := audit.Add(event)

	if err != nil {
		t.Fatalf("Add returned error: %s", err)
		return
	}

	givenMessages, err := audit.GetEvents(timestamp)

	if err != nil {
		t.Fatalf("GetEvents returned error: %s", err)
		return
	}

	if len(givenMessages) == 0 {
		t.Fatalf("Given Messages got length of 0")
		return
	}

	givenMessage := givenMessages[0].(CaseTriggerAuditEvent)
	if !reflect.DeepEqual(givenMessage, event) {
		t.Fatalf("Given %+v, expected %+v", givenMessage, messages[0])
	}
}
