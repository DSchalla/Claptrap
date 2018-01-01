package claptrap

type Config struct {
	BotToken   string
	AdminToken string
	CaseDir    string
}

func NewConfig(botToken, adminToken, configDir string) Config {
	config := Config{
		BotToken:   botToken,
		AdminToken: adminToken,
		CaseDir:    configDir,
	}
	return config
}
