package rules

import (
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"reflect"
	"encoding/gob"
	"testing"
	"bytes"
)

func TestCaseManager_GetForType(t *testing.T) {
	var buffer bytes.Buffer

	cases := []Case{
		{
			Name: "TestCase",
			Intercept: true,
			Conditions: []Condition{
				TextEqualsCondition{
					Condition: "Hello World",
				},
			},
			Responses: []Response{
				MessageChannelResponse{
					Message: "Hello User",
				},
			},
		},
		{
			Name: "TestCase #2",
			Intercept: false,
			Conditions: []Condition{
				TextStartsWithCondition{
					Condition: "!greet",
				},
				UserIsRoleCondition{
					Condition: "system_admin",
				},
			},
			Responses: []Response{
				MessageUserResponse{
					Message: "Hello User",
				},
				DeleteMessageResponse{},
			},
		},
	}

	gob.Register(TextEqualsCondition{})
	gob.Register(TextStartsWithCondition{})
	gob.Register(UserIsRoleCondition{})
	gob.Register(MessageUserResponse{})
	gob.Register(MessageChannelResponse{})
	gob.Register(DeleteMessageResponse{})

	enc := gob.NewEncoder(&buffer)
	enc.Encode(cases)

	api := &plugintest.API{}
	api.On("KVGet", "cases.message").Return(buffer.Bytes(), nil)

	audit := NewCaseManager(api)
	givenCases, err := audit.GetForType("message")

	if err != nil {
		t.Fatalf("GetEvents returned error: %s", err)
		return
	}

	if len(givenCases) == 0 {
		t.Fatalf("Given Messages got length of 0")
		return
	}

	givenCase := givenCases[0]
	if !reflect.DeepEqual(givenCase, cases[0]) {
		t.Fatalf("Given %+v, expected %+v", givenCase, cases[0])
	}

	givenCase = givenCases[1]
	if !reflect.DeepEqual(givenCase, cases[1]) {
		t.Fatalf("Given %+v, expected %+v", givenCase, cases[1])
	}
}

func TestCaseManager_Add(t *testing.T) {
	var buffer bytes.Buffer

	newCase := Case{
		Name: "TestCase",
		Intercept: true,
		Conditions: []Condition{
			TextEqualsCondition{
				Condition: "Hello World",
			},
		},
		Responses: []Response{
			MessageChannelResponse{
				Message: "Hello User",
			},
		},
	}

	cases := []Case{
		newCase,
	}

	gob.Register(TextEqualsCondition{})
	gob.Register(MessageChannelResponse{})

	enc := gob.NewEncoder(&buffer)
	enc.Encode(cases)

	api := &plugintest.API{}
	api.On("KVSet", "cases.message", buffer.Bytes()).Return(nil)
	api.On("KVGet", "cases.message").Return(nil, nil)

	cm := NewCaseManager(api)
	err := cm.Add("message", newCase)

	if err != nil {
		t.Fatalf("Add returned error: %s", err)
		return
	}

}
