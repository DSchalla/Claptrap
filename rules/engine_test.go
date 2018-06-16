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

	e := rules.NewEngine("")
	e.AddCase("message", testCase)

	event := provider.Event{
		UserID:      "UABCDEF",
		UserName:    "ABCDEF",
		ChannelID:   "CABCDEF",
		ChannelName: "general",
		Text:        "abcdef",
	}
	hit := e.EvaluateEvent(event)

	if hit {
		t.Errorf("Expected True, got False")
	}

	event = provider.Event{
		UserID:      "UABCDEF",
		UserName:    "ABCDEF",
		ChannelID:   "CABCDEF",
		ChannelName: "general",
		Text:        "foobar",
	}
	hit = e.EvaluateEvent(event)

	if hit {
		t.Errorf("Expected False, got True")
	}
}

func TestEngine_ComplexCondition(t *testing.T) {
	e := rules.NewEngine("")
	e.Com
}