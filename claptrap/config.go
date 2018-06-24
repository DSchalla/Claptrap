package claptrap

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Name 	   string
	Cases 		string
	AutoJoinAllChannel bool
}

func NewConfig(filePath string) Config {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("[!] Error reading config file")
	}
	config := &Config{}
	yaml.Unmarshal(content, config)
	return *config
}
