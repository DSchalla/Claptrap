package rules_test

import (
	"github.com/DSchalla/Claptrap/rules"
	"testing"
	"github.com/DSchalla/Claptrap/provider"
)

func TestTextContainsCondition(t *testing.T) {
	cond := rules.TextContainsCondition{
		Condition: "Test",
	}
	if !cond.Test(provider.Event{Text: "Test 123"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{Text: "test 123"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestTextEqualsCondition(t *testing.T) {
	cond := rules.TextEqualsCondition{
		Condition: "Foobar",
	}

	if !cond.Test(provider.Event{Text: "Foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{Text: "FooBar"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestTextStartsWithCondition(t *testing.T) {
	cond := rules.TextStartsWithCondition{
		Condition: "Foobar",
	}

	if !cond.Test(provider.Event{Text: "Foobar abc"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{Text: "abc Foobar"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestUserEqualsCondition_Test(t *testing.T) {
	cond := rules.UserEqualsCondition{
		Condition: "foobar",
		Parameter: "user",
	}

	if !cond.Test(provider.Event{UserName: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if !cond.Test(provider.Event{UserID: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{UserID: "abcdef"}) {
		t.Errorf("Expected False, got True")
	}

	cond = rules.UserEqualsCondition{
		Condition: "foobar",
		Parameter: "actor",
	}

	if !cond.Test(provider.Event{ActorName: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if !cond.Test(provider.Event{ActorID: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{ActorID: "abcdef"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestUserIsRoleCondition_Test(t *testing.T) {
	cond := rules.UserIsRoleCondition{
		Condition: "admin",
		Parameter: "user",
	}

	if !cond.Test(provider.Event{UserRole: "admin"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{UserRole: "user"}) {
		t.Errorf("Expected False, got True")
	}

	cond = rules.UserIsRoleCondition{
		Condition: "admin",
		Parameter: "actor",
	}

	if !cond.Test(provider.Event{ActorRole: "admin"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{ActorRole: "user"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestChannelEqualsCondition_Test(t *testing.T) {
	cond := rules.ChannelEqualsCondition{
		Condition: "foobar",
	}

	if !cond.Test(provider.Event{ChannelID: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if !cond.Test(provider.Event{ChannelName: "foobar"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(provider.Event{ChannelID: "abcdef"}) {
		t.Errorf("Expected False, got True")
	}

	if cond.Test(provider.Event{ChannelName: "abcdef"}) {
		t.Errorf("Expected False, Got True")
	}
}

func TestChannelIsTypeCondition_Test(t *testing.T) {
	cond := rules.ChannelIsTypeCondition{
		Condition: "channel",
	}

	if !cond.Test(provider.Event{ChannelID: "C123456"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(provider.Event{ChannelID: "D123546"}) {
		t.Errorf("Expected False, got True")
	}

	cond = rules.ChannelIsTypeCondition{
		Condition: "group",
	}

	if !cond.Test(provider.Event{ChannelID: "G123456"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(provider.Event{ChannelID: "C123546"}) {
		t.Errorf("Expected False, got True")
	}

	cond = rules.ChannelIsTypeCondition{
		Condition: "dm",
	}

	if !cond.Test(provider.Event{ChannelID: "D123456"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(provider.Event{ChannelID: "C123546"}) {
		t.Errorf("Expected False, got True")
	}
}
