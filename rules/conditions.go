package rules

import (
	"log"
	"math/rand"
	"strings"
	"github.com/DSchalla/Claptrap/provider"
	"regexp"
)

type Condition interface {
	Test(event provider.Event) bool
}

type TextContainsCondition struct {
	Condition string
}

func (t TextContainsCondition) Test(event provider.Event) bool {
	return strings.Contains(event.Text, t.Condition)
}

type TextEqualsCondition struct {
	Condition string
}

func (t TextEqualsCondition) Test(event provider.Event) bool {
	return event.Text == t.Condition
}

type TextStartsWithCondition struct {
	Condition string
}

func (t TextStartsWithCondition) Test(event provider.Event) bool {
	return strings.HasPrefix(event.Text, t.Condition)
}

func NewTextMatchesCondition(expression string) (*TextMatchesCondition, error){
	var err error
	c := &TextMatchesCondition{}
	c.expression = expression
	c.regexp, err = regexp.Compile(expression)
	if err != nil {
		return nil, err
	}
	return c, err
}

type TextMatchesCondition struct {
	expression string
	regexp *regexp.Regexp
}

func (t TextMatchesCondition) Test(event provider.Event) bool {
	return t.regexp.MatchString(event.Text)
}

type RandomCondition struct {
	Likeness int
}

func (t RandomCondition) Test(event provider.Event) bool {
	randomInt := rand.Int()
	return t.Likeness > (randomInt % 100)
}

type UserEqualsCondition struct {
	Condition string
	Parameter string
}

func (u UserEqualsCondition) Test(event provider.Event) bool {
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
		return false
	}

	return (userID == u.Condition) || (userName == u.Condition)
}

type UserIsRoleCondition struct {
	Condition string
	Parameter string
}

func (u UserIsRoleCondition) Test(event provider.Event) bool {
	role := ""

	if u.Parameter == "" || u.Parameter == "user" {
		role = event.UserRole
	} else if u.Parameter == "actor" {
		role = event.ActorRole
	} else {
		log.Printf("[!] Error: Unknown Parameter for UserIsRoleCondition: '%s' \n", u.Parameter)
		return false
	}

	return strings.Contains(role, u.Condition)
}

type ChannelEqualsCondition struct {
	Condition string
}

func (c ChannelEqualsCondition) Test(event provider.Event) bool {
	return (event.ChannelID == c.Condition) || (event.ChannelName == c.Condition)
}

type ChannelIsTypeCondition struct {
	Condition string
}

func (c ChannelIsTypeCondition) Test(event provider.Event) bool {
	condition := ""

	if c.Condition == "public" {
		condition = "O"
	} else if c.Condition == "private" {
		condition = "P"
	} else {
		condition = "D"
	}

	return event.ChannelType == condition
}
