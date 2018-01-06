package claptrap

type Config struct {
	BotToken   string
	AdminToken string
	CaseDir    string
	AutoJoinAllChannel bool
}

func NewConfig(botToken, adminToken, configDir string, autoJoinAllChannel bool) Config {
	config := Config{
		BotToken:   botToken,
		AdminToken: adminToken,
		CaseDir:    configDir,
		AutoJoinAllChannel: autoJoinAllChannel,
	}
	return config
}
