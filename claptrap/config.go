package claptrap

type Config struct {
	ApiUrl   string
	Username string
	Password string
	Team string
	CaseDir    string
	AutoJoinAllChannel bool
}

func NewConfig(apiUrl, username, password, team, configDir string, autoJoinAllChannel bool) Config {
	config := Config{
		ApiUrl: apiUrl,
		Username:   username,
		Password: password,
		Team: team,
		CaseDir:    configDir,
		AutoJoinAllChannel: autoJoinAllChannel,
	}
	return config
}
