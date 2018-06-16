package rules_test

import (
	"github.com/DSchalla/Claptrap/provider"
	"github.com/DSchalla/Claptrap/rules"
	"testing"
)

func TestTextContainsCondition(t *testing.T) {
	cond, _ := rules.NewTextContainsCondition("Test")
	if !cond.Test(provider.Event{Text: "Test 123"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{Text: "test 123"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestTextEqualsCondition(t *testing.T) {
	cond, _ := rules.NewTextEqualsCondition("Foobar")

	if !cond.Test(provider.Event{Text: "Foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{Text: "FooBar"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestTextStartsWithCondition(t *testing.T) {
	cond, _ := rules.NewTextStartsWithCondition("Foobar")

	if !cond.Test(provider.Event{Text: "Foobar abc"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{Text: "abc Foobar"}) {
		t.Errorf("Expected False, got True")
	}
}

func TestTextMatchesCondition(t *testing.T) {
	cond, err := rules.NewTextMatchesCondition("^a[0-9]b$")

	if err != nil {
		t.Errorf("Expected nil, got error")
	}

	if !cond.Test(provider.Event{Text: "a3b"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{Text: "abc a3b"}) {
		t.Errorf("Expected False, got True")
	}

	if cond.Test(provider.Event{Text: "a3b abc"}) {
		t.Errorf("Expected False, got True")
	}

	_, err = rules.NewTextMatchesCondition("^a([0-9]]))]b$")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUserEqualsCondition_Test(t *testing.T) {
	cond, _ := rules.NewUserEqualsCondition("foobar", "user")

	if !cond.Test(provider.Event{UserName: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if !cond.Test(provider.Event{UserID: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{UserID: "abcdef"}) {
		t.Errorf("Expected False, got True")
	}

	cond, _ = rules.NewUserEqualsCondition("foobar", "actor")

	if !cond.Test(provider.Event{ActorName: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if !cond.Test(provider.Event{ActorID: "foobar"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{ActorID: "abcdef"}) {
		t.Errorf("Expected False, got True")
	}

	_, err := rules.NewUserEqualsCondition("foobar", "doesntexist")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUserIsRoleCondition_Test(t *testing.T) {
	cond, _ := rules.NewUserIsRoleCondition("admin", "user")

	if !cond.Test(provider.Event{UserRole: "admin"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{UserRole: "user"}) {
		t.Errorf("Expected False, got True")
	}

	cond, _ = rules.NewUserIsRoleCondition("admin", "actor")

	if !cond.Test(provider.Event{ActorRole: "admin"}) {
		t.Errorf("Expected True, got False")
	}

	if cond.Test(provider.Event{ActorRole: "user"}) {
		t.Errorf("Expected False, got True")
	}

	_, err := rules.NewUserIsRoleCondition("foobar", "doesntexist")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestChannelEqualsCondition_Test(t *testing.T) {
	cond, _ := rules.NewChannelEqualsCondition("foobar")

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
	cond, _ := rules.NewChannelIsTypeCondition("public")

	if !cond.Test(provider.Event{ChannelType: "O"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(provider.Event{ChannelType: "D"}) {
		t.Errorf("Expected False, got True")
	}

	cond, _ = rules.NewChannelIsTypeCondition("private")

	if !cond.Test(provider.Event{ChannelType: "P"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(provider.Event{ChannelType: "D"}) {
		t.Errorf("Expected False, got True")
	}

	cond, _ = rules.NewChannelIsTypeCondition("dm")

	if !cond.Test(provider.Event{ChannelType: "D"}) {
		t.Errorf("Expected True, Got False")
	}

	if cond.Test(provider.Event{ChannelType: "O"}) {
		t.Errorf("Expected False, got True")
	}
}
