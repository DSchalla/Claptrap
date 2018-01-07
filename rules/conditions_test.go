package rules_test

import (
	"github.com/DSchalla/Claptrap/rules"
	"testing"
)

func TestTextContainsCondition(t *testing.T) {
	cond := rules.TextContainsCondition{
		Condition: "Test",
	}
	if !cond.Test(rules.Event{Text: "Test 123"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(rules.Event{Text: "test 123"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestTextEqualsCondition(t *testing.T) {
	cond := rules.TextEqualsCondition{
		Condition: "Foobar",
	}

	if !cond.Test(rules.Event{Text: "Foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(rules.Event{Text: "FooBar"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestTextStartsWithCondition(t *testing.T) {
	cond := rules.TextStartsWithCondition{
		Condition: "Foobar",
	}

	if !cond.Test(rules.Event{Text: "Foobar abc"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(rules.Event{Text: "abc Foobar"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestUserEqualsCondition_Test(t *testing.T) {
	cond := rules.UserEqualsCondition{
		Condition: "foobar",
		Parameter: "user",
	}

	if !cond.Test(rules.Event{UserName: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if !cond.Test(rules.Event{UserID: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(rules.Event{UserID: "abcdef"}) {
		t.Errorf("Expected False, got True")
	}

	cond = rules.UserEqualsCondition{
		Condition: "foobar",
		Parameter: "inviter",
	}

	if !cond.Test(rules.Event{InviterName: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if !cond.Test(rules.Event{InviterID: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(rules.Event{InviterID: "abcdef"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestUserIsRoleCondition_Test(t *testing.T) {
	cond := rules.UserIsRoleCondition{
		Condition: "admin",
		Parameter: "user",
	}

	if !cond.Test(rules.Event{UserRole: "admin"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(rules.Event{UserRole: "user"}) {
		t.Errorf("Expected False, got True")
	}

	cond = rules.UserIsRoleCondition{
		Condition: "admin",
		Parameter: "inviter",
	}

	if !cond.Test(rules.Event{InviterRole: "admin"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(rules.Event{InviterRole: "user"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestChannelEqualsCondition_Test(t *testing.T) {
	cond := rules.ChannelEqualsCondition{
		Condition: "foobar",
	}

	if !cond.Test(rules.Event{ChannelID: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if !cond.Test(rules.Event{ChannelName: "foobar"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(rules.Event{ChannelID: "abcdef"}) {
		t.Errorf("Expected False, got True")
	}

	if cond.Test(rules.Event{ChannelName: "abcdef"}) {
		t.Errorf("Expected False, Got True")
	}
}

func TestChannelIsTypeCondition_Test(t *testing.T) {
	cond := rules.ChannelIsTypeCondition{
		Condition: "channel",
	}

	if !cond.Test(rules.Event{ChannelID: "C123456"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(rules.Event{ChannelID: "D123546"}) {
		t.Errorf("Expected False, got True")
	}

	cond = rules.ChannelIsTypeCondition{
		Condition: "group",
	}

	if !cond.Test(rules.Event{ChannelID: "G123456"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(rules.Event{ChannelID: "C123546"}) {
		t.Errorf("Expected False, got True")
	}

	cond = rules.ChannelIsTypeCondition{
		Condition: "dm",
	}

	if !cond.Test(rules.Event{ChannelID: "D123456"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(rules.Event{ChannelID: "C123546"}) {
		t.Errorf("Expected False, got True")
	}
}
