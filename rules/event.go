package rules

type Event struct {
	Type        string
	UserID      string
	UserName    string
	UserRole    string
	InviterID   string
	InviterName string
	InviterRole string
	ChannelID   string
	ChannelName string
	Text        string
	Timestamp   string
}
