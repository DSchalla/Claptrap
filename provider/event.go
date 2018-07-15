package provider

type Event struct {
	Type        string
	PostID      string
	UserID      string
	UserName    string
	UserRole    string
	ActorID     string
	ActorName   string
	ActorRole   string
	ChannelID   string
	ChannelName string
	ChannelType string
	TeamID      string
	TeamName 	string
	Text        string
	Timestamp   int64
}
