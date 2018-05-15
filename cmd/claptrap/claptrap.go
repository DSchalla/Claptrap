package main

import (
	"flag"
	"github.com/DSchalla/Claptrap/claptrap"
	"log"
)

func main() {
	var configFile = flag.String("config_file", "config.yaml", "config.yaml")
	flag.Parse()

	config := claptrap.NewConfig(*configFile)
	botServer, err := claptrap.NewBotServer(config)
	if err != nil {
		log.Fatalf("Error starting botserver: %s\n", err)
	}
	botServer.Start()
}
