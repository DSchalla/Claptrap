package claptrap

type Config struct {
	BotToken   string
	AdminToken string
	ConfigDir  string
}

func NewConfig(botToken, adminToken, configDir string) Config {
	config := Config{
		BotToken:   botToken,
		AdminToken: adminToken,
		ConfigDir:  configDir,
	}
	return config
}
