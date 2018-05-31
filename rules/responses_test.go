package rules_test

import (
	"testing"
	"github.com/DSchalla/Claptrap/rules"
	"os"
	"github.com/DSchalla/Claptrap/provider"
)

var prov *provider.Debug
var event provider.Event

func TestMain(m *testing.M) {
	prov = provider.NewDebug()
	event = provider.Event{
		PostID:		 "Post1",
		UserID:      "User1",
		UserName:    "Username1",
		ActorID: 	 "Actor1",
		ActorName: 	 "Actorname1",
		ChannelID:   "Channel1",
		ChannelName: "Channelname1",
		Text:        "hunter2",
	}
	retCode := m.Run()
	os.Exit(retCode)
}

func TestMessageChannelResponse(t *testing.T) {
	resp, err := rules.NewMessageChannelResponse("", "Hello World")
	res := resp.Execute(prov, event)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if res != true {
		t.Errorf("Expected true, got false")
	}

	if prov.MessagePublicLog[0].Message != "Hello World" {
		t.Errorf("Expected Hello World, got %s", prov.MessagePublicLog[1].Message)
	}

	if prov.MessagePublicLog[0].ChannelID != "Channel1" {
		t.Errorf("Expected Channel1, got %s", prov.MessagePublicLog[1].ChannelID)
	}

	resp, _ = rules.NewMessageChannelResponse("Channel2", "Hello World")
	res = resp.Execute(prov, event)

	if res != true {
		t.Errorf("Expected true, got false")
	}

	if len(prov.MessagePublicLog) < 2 {
		t.Errorf("Expected len == 2 of message log, got < 2")
	}

	if prov.MessagePublicLog[1].ChannelID != "Channel2" {
		t.Errorf("Expected Channel2, got %s", prov.MessagePublicLog[1].ChannelID)
	}
}

func TestMessageUserResponse(t *testing.T) {
	resp, err := rules.NewMessageUserResponse("", "Hello World")
	res := resp.Execute(prov, event)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if res != true {
		t.Errorf("Expected true, got false")
	}

	if prov.MessageUserLog[0].Message != "Hello World" {
		t.Errorf("Expected Hello World, got %s", prov.MessageUserLog[1].Message)
	}

	if prov.MessageUserLog[0].UserID != "User1" {
		t.Errorf("Expected Channel1, got %s", prov.MessageUserLog[1].UserID)
	}

	resp, _ = rules.NewMessageUserResponse("User2", "Hello World")
	res = resp.Execute(prov, event)

	if res != true {
		t.Errorf("Expected true, got false")
	}

	if len(prov.MessageUserLog) < 2 {
		t.Errorf("Expected len == 2 of message log, got < 2")
	}

	if prov.MessageUserLog[1].UserID != "User2" {
		t.Errorf("Expected User2, got %s", prov.MessageUserLog[1].UserID)
	}
}

func TestInviteUserResponse(t *testing.T) {
	resp, err := rules.NewInviteUserResponse("Channel1", "")
	res := resp.Execute(prov, event)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if res != true {
		t.Errorf("Expected true, got false")
	}

	if prov.InviteUserLog[0].ChannelID != "Channel1" {
		t.Errorf("Expected Channel1, got %s", prov.InviteUserLog[1].ChannelID)
	}

	if prov.InviteUserLog[0].UserID != "User1" {
		t.Errorf("Expected User1, got %s", prov.InviteUserLog[1].UserID)
	}

	resp, _ = rules.NewInviteUserResponse("Channel2", "User2")
	res = resp.Execute(prov, event)

	if res != true {
		t.Errorf("Expected true, got false")
	}

	if len(prov.InviteUserLog) < 2 {
		t.Errorf("Expected len == 2 of message log, got < 2")
	}

	if prov.InviteUserLog[1].UserID != "User2" {
		t.Errorf("Expected User2, got %s", prov.InviteUserLog[1].UserID)
	}
}


func TestKickUserResponse(t *testing.T) {
	resp, err := rules.NewKickUserResponse("", "")
	res := resp.Execute(prov, event)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if res != true {
		t.Errorf("Expected true, got false")
	}

	if prov.KickUserLog[0].ChannelID != "Channel1" {
		t.Errorf("Expected Channel1, got %s", prov.KickUserLog[1].ChannelID)
	}

	if prov.KickUserLog[0].UserID != "User1" {
		t.Errorf("Expected User1, got %s", prov.KickUserLog[1].UserID)
	}

	resp, _ = rules.NewKickUserResponse("Channel2", "User2")
	res = resp.Execute(prov, event)

	if res != true {
		t.Errorf("Expected true, got false")
	}

	if len(prov.KickUserLog) < 2 {
		t.Errorf("Expected len == 2 of message log, got < 2")
	}

	if prov.KickUserLog[1].ChannelID != "Channel2" {
		t.Errorf("Expected Channel2, got %s", prov.KickUserLog[1].ChannelID)
	}

	if prov.KickUserLog[1].UserID != "User2" {
		t.Errorf("Expected User2, got %s", prov.KickUserLog[1].UserID)
	}
}


func TestDeleteMessageResponse(t *testing.T) {
	resp, err := rules.NewDeleteMessageResponse()
	res := resp.Execute(prov, event)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if res != true {
		t.Errorf("Expected true, got false")
	}

	if prov.DeleteMessageLog[0].PostID != "Post1" {
		t.Errorf("Expected Post1, got %s", prov.DeleteMessageLog[1].PostID)
	}
}
