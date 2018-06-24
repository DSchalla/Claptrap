package rules_test

import (
	"github.com/DSchalla/Claptrap/provider"
	"github.com/DSchalla/Claptrap/rules"
	"testing"
)

func TestEngine_EvaluateMessageEvent(t *testing.T) {
	testCase := rules.Case{
		Name: "Example Case",
		Conditions: []rules.Condition{
			rules.TextContainsCondition{
				Condition: "abc",
			},
		},
		Responses: nil,
	}

	e := rules.NewEngine(nil)
	e.AddCase("message", testCase)

	event := provider.Event{
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