package rules

import (
	"log"
	"math/rand"
	"strings"
)

type Condition interface {
	Test(event Event) bool
}

type TextContainsCondition struct {
	Condition string
}

func (t TextContainsCondition) Test(event Event) bool {
	return strings.Contains(event.Text, t.Condition)
}

type TextEqualsCondition struct {
	Condition string
}

func (t TextEqualsCondition) Test(event Event) bool {
	return event.Text == t.Condition
}

type TextStartsWithCondition struct {
	Condition string
}

func (t TextStartsWithCondition) Test(event Event) bool {
	return strings.HasPrefix(event.Text, t.Condition)
}

type RandomCondition struct {
	Likeness int
}

func (t RandomCondition) Test(event Event) bool {
	randomInt := rand.Int()
	return t.Likeness > (randomInt % 100)
}

type UserEqualsCondition struct {
	Condition string
	Parameter string
}

func (u UserEqualsCondition) Test(event Event) bool {
	userID := ""
	userName := ""

	if u.Parameter == "" || u.Parameter == "user" {
		userID = event.UserID
		userName = event.UserName
	} else if u.Parameter == "actor" {
		userID = event.ActorID
		userName = event.ActorName
	} else {
		log.Printf("[!] Error: Unknown Parameter for UserIDEqaulsCondition: '%s' \n", u.Parameter)
	}

	return (userID == u.Condition) || (userName == u.Condition)
}

type UserIsRoleCondition struct {
	Condition string
	Parameter string
}

func (u UserIsRoleCondition) Test(event Event) bool {
	role := ""

	if u.Parameter == "" || u.Parameter == "user" {
		role = event.UserRole
	} else if u.Parameter == "actor" {
		role = event.ActorRole
	} else {
		log.Printf("[!] Error: Unknown Parameter for UserIsRoleCondition: '%s' \n", u.Parameter)
	}

	return strings.Contains(role, u.Condition)
}

type ChannelEqualsCondition struct {
	Condition string
}

func (c ChannelEqualsCondition) Test(event Event) bool {
	return (event.ChannelID == c.Condition) || (event.ChannelName == c.Condition)
}

type ChannelIsTypeCondition struct {
	Condition string
}

func (c ChannelIsTypeCondition) Test(event Event) bool {
	prefix := ""

	if c.Condition == "channel" {
		prefix = "C"
	} else if c.Condition == "group" {
		prefix = "G"
	} else {
		prefix = "D"
	}

	return strings.HasPrefix(event.ChannelID, prefix)
}
