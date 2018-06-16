package claptrap

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	General    GeneralConfig    `yaml:"general"`
	Mattermost MattermostConfig `yaml:"mattermost"`
	Webserver  WebserverConfig  `yaml:"webserver"`
}

type GeneralConfig struct {
	CaseDir            string `yaml:"case_dir"`
	AutoJoinAllChannel bool   `yaml:"auto_join_all_channel"`
}

type MattermostConfig struct {
	ApiUrl   string `yaml:"api_url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Team     string `yaml:"team"`
}

type WebserverConfig struct {
	Enabled bool   `yaml:"enabled"`
	Listen  string `yaml:"listen"`
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
