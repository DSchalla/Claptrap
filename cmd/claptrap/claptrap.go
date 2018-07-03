package main

import (
	"net/http"

	"github.com/mattermost/mattermost-server/model"

	"github.com/DSchalla/Claptrap/claptrap"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/plugin"
	"fmt"
	"errors"
)

type ClaptrapPlugin struct {
	plugin.MattermostPlugin
	claptrap *claptrap.BotServer
	config   *claptrap.Config
}

func (c *ClaptrapPlugin) OnActivate() error {
	mlog.Debug("[CLAPTRAP-PLUGIN] OnActivate Hook Start")
	var err error
	c.readConfig()
	c.claptrap, err = claptrap.NewBotServer(c.API, *c.config)
	if err != nil {
		mlog.Debug(fmt.Sprintf("[CLAPTRAP-PLUGIN]  NewBotServer returned error: %s", err))
	}
	mlog.Debug("[CLAPTRAP-PLUGIN] OnActivate Hook End")

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (c *ClaptrapPlugin) OnDeactivate() error {
	if c.claptrap != nil {
		c.claptrap.Shutdown()
	}

	return nil
}

func (c *ClaptrapPlugin) OnConfigurationChange() error {
	err := c.readConfig()
	if err != nil {
		return err
	}
	return c.reloadConfig()
}

func (c *ClaptrapPlugin) MessageWillBePosted(post *model.Post) (*model.Post, string) {
	if post.Props["from_claptrap"] != nil && post.Props["from_claptrap"].(bool) == true {
		return post, ""
	}
	mlog.Debug("[CLAPTRAP-PLUGIN] MessageWillBePosted Hook Start")
	post, rejectMessage := c.claptrap.HandleMessage(post, true)
	mlog.Debug("[CLAPTRAP-PLUGIN] MessageWillBePosted Hook End")
	return post, rejectMessage
}

func (c *ClaptrapPlugin) MessageHasBeenPosted(post *model.Post) {
	if post.Props["from_claptrap"] != nil && post.Props["from_claptrap"].(bool) == true {
		return
	}
	mlog.Debug("[CLAPTRAP-PLUGIN] MessageHasBeenPosted Hook Start")
	c.claptrap.HandleMessage(post, false)
	mlog.Debug("[CLAPTRAP-PLUGIN] MessageHasBeenPosted Hook End")
}

func (c *ClaptrapPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mlog.Debug("[CLAPTRAP-PLUGIN] ServeHTTP Hook Start")
	c.claptrap.HandleHTTP(w, r)
	mlog.Debug("[CLAPTRAP-PLUGIN] ServeHTTP Hook End")
}

func (c *ClaptrapPlugin) readConfig() error {
	c.config = &claptrap.Config{}
	err := c.API.LoadPluginConfiguration(c.config)
	return err
}

func (c *ClaptrapPlugin) reloadConfig() error {
	if c.claptrap != nil {
		c.claptrap.ReloadConfig(*c.config)
	}

	return nil
}

func main() {
	plugin.ClientMain(&ClaptrapPlugin{})
}
