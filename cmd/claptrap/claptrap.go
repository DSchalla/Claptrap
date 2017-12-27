package main

import (
	"github.com/DSchalla/Claptrap/claptrap"
	"sync"
	"flag"
)

func main() {
	var botToken = flag.String("bot_token", "", "xoxb-0x0x0x0x0x0x0x0x0x0x")
	var adminToken = flag.String("admin_token", "", "xoxs-0x0x0x0x0x0x0x0x0x0x")
	var configDir = flag.String("config_dir", "config/", "config/")
	flag.Parse()

	config := claptrap.NewConfig(*botToken, *adminToken, *configDir)
	botServer := claptrap.NewBotServer(config)
	botServer.Start()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
