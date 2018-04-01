package rules

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
	Text        string
	Timestamp   int64
}
