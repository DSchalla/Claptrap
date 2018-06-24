package main

import (
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/plugin/rpcplugin"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/model"

	"github.com/DSchalla/Claptrap/claptrap"
	"github.com/mattermost/mattermost-server/mlog"
)

type ClaptrapPlugin struct{
	api plugin.API
	claptrap *claptrap.BotServer
	botuser  *model.User
	config   *claptrap.Config
}

func (c *ClaptrapPlugin) OnActivate(api plugin.API) error {
	mlog.Debug("[CLAPTRAP-PLUGIN] OnActivate Hook Start")
	var err error
	c.api = api
	c.readConfig()
	c.claptrap, err = claptrap.NewBotServer(api, *c.config)
	mlog.Debug("[CLAPTRAP-PLUGIN] OnActivate Hook End")
	return err
}

func (c *ClaptrapPlugin) OnConfigurationChange() error {
	err := c.readConfig()
	if err != nil {
		return err
	}
	return c.reloadConfig()
}

func (c *ClaptrapPlugin) MessageWillBePosted(post *model.Post) (*model.Post, string){
	mlog.Debug("[CLAPTRAP-PLUGIN] MessageWillBePosted Hook Start")
	post, rejectMessage := c.claptrap.HandleMessage(post, true)
	mlog.Debug("[CLAPTRAP-PLUGIN] MessageWillBePosted Hook End")
	return post, rejectMessage
}

func (c *ClaptrapPlugin) MessageHasBeenPosted(post *model.Post) {
	mlog.Debug("[CLAPTRAP-PLUGIN] MessageHasBeenPosted Hook Start")
	c.claptrap.HandleMessage(post, false)
	mlog.Debug("[CLAPTRAP-PLUGIN] MessageHasBeenPosted Hook End")
}

func (c *ClaptrapPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func (c *ClaptrapPlugin) readConfig() error {
	c.config = &claptrap.Config{}
	err := c.api.LoadPluginConfiguration(c.config)
	c.config.Name = "claptrap"
	c.config.Cases = `[{  "name": "Regexp Message",  "conditions": [    {"type": "text_matches", "condition": "^a[0-9]b$"}  ],  "responses": [    {"action": "message_channel", "message": "Yes, Regexp works!"}  ]}]`
	return err
}

func (c *ClaptrapPlugin) reloadConfig() error {
	if c.claptrap != nil {
		c.claptrap.ReloadConfig(*c.config)
	}
	return nil
}

func main() {
	rpcplugin.Main(&ClaptrapPlugin{})
}