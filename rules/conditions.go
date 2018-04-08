package rules

import (
	"math/rand"
	"strings"
	"github.com/DSchalla/Claptrap/provider"
	"regexp"
	"fmt"
)

type Condition interface {
	Test(event provider.Event) bool
}

func NewTextContainsCondition(condition string) (*TextContainsCondition, error){
	return &TextContainsCondition{Condition: condition}, nil
}

type TextContainsCondition struct {
	Condition string
}

func (t TextContainsCondition) Test(event provider.Event) bool {
	return strings.Contains(event.Text, t.Condition)
}

func NewTextEqualsCondition(condition string) (*TextEqualsCondition, error){
	return &TextEqualsCondition{Condition: condition}, nil
}

type TextEqualsCondition struct {
	Condition string
}

func (t TextEqualsCondition) Test(event provider.Event) bool {
	return event.Text == t.Condition
}

func NewTextStartsWithCondition(condition string) (*TextStartsWithCondition, error){
	return &TextStartsWithCondition{Condition: condition}, nil
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

func NewRandomCondition(likeness int) (*RandomCondition, error){
	return &RandomCondition{Likeness: likeness}, nil
}

type RandomCondition struct {
	Likeness int
}

func (t RandomCondition) Test(event provider.Event) bool {
	randomInt := rand.Int()
	return t.Likeness > (randomInt % 100)
}

func NewUserEqualsCondition(condition, parameter string) (*UserEqualsCondition, error){

	if parameter != "" && parameter != "user" && parameter != "actor" {
		return nil, fmt.Errorf("unknown Parameter for UserIDEqualsCondition: '%s'", parameter)
	}

	return &UserEqualsCondition{Condition: condition, Parameter: parameter}, nil
}

type UserEqualsCondition struct {
	Condition string
	Parameter string
}

func (u UserEqualsCondition) Test(event provider.Event) bool {
	userID := ""
	userName := ""

	if u.Parameter == "actor" {
		userID = event.UserID
		userName = event.UserName
		userID = event.ActorID
		userName = event.ActorName
	} else {
		userID = event.UserID
		userName = event.UserName
	}

	return (userID == u.Condition) || (userName == u.Condition)
}

func NewUserIsRoleCondition(condition, parameter string) (*UserIsRoleCondition, error){

	if parameter != "" && parameter != "user" && parameter != "actor" {
		return nil, fmt.Errorf("unknown Parameter for UserIsRoleCondition: '%s'", parameter)
	}

	return &UserIsRoleCondition{Condition: condition, Parameter: parameter}, nil
}

type UserIsRoleCondition struct {
	Condition string
	Parameter string
}

func (u UserIsRoleCondition) Test(event provider.Event) bool {
	role := ""

	if u.Parameter == "actor" {
		role = event.ActorRole
	} else {
		role = event.UserRole
	}

	return strings.Contains(role, u.Condition)
}

func NewChannelEqualsCondition(condition string) (*ChannelEqualsCondition, error){
	return &ChannelEqualsCondition{Condition: condition}, nil
}

type ChannelEqualsCondition struct {
	Condition string
}

func (c ChannelEqualsCondition) Test(event provider.Event) bool {
	return (event.ChannelID == c.Condition) || (event.ChannelName == c.Condition)
}


func NewChannelIsTypeCondition(condition string) (*ChannelIsTypeCondition, error){
	return &ChannelIsTypeCondition{Condition: condition}, nil
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
