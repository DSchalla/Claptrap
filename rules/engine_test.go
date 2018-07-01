package rules

import (
	"github.com/DSchalla/Claptrap/provider"
	"testing"
	"encoding/gob"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"bytes"
)

func TestEngine_EvaluateMessageEvent(t *testing.T) {
	var buffer bytes.Buffer

	testCase := Case{
		Name: "Example Case",
		Conditions: []Condition{
			TextContainsCondition{
				Condition: "abc",
			},
		},
		Responses: nil,
	}

	gob.Register(TextContainsCondition{})

	enc := gob.NewEncoder(&buffer)
	enc.Encode([]Case{testCase})

	api := &plugintest.API{}
	api.On("KVSet", "cases.message", buffer.Bytes()).Return(nil)
	api.On("KVGet", "cases.message").Return(nil, nil)


	cm := NewCaseManager(api)
	cm.Add("message", testCase)

	api.On("KVGet", "cases.message").Return(buffer.Bytes(), nil)

	e := NewEngine(cm, nil, nil)

	event := provider.Event{
		Type:		 "message",
		UserID:      "UABCDEF",
		UserName:    "ABCDEF",
		ChannelID:   "CABCDEF",
		ChannelName: "general",
		Text:        "abcdef",
	}
	hit := e.EvaluateEvent(event, false)

	if hit.Hit {
		t.Errorf("Expected True, got False")
	}

	event = provider.Event{
		Type:		 "message",
		UserID:      "UABCDEF",
		UserName:    "ABCDEF",
		ChannelID:   "CABCDEF",
		ChannelName: "general",
		Text:        "foobar",
	}
	hit = e.EvaluateEvent(event, false)

	if hit.Hit {
		t.Errorf("Expected False, got True")
	}
}
