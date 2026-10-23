package common

// RegisterArgs is used when a user attempts to join
type RegisterArgs struct {
	Username string
}

// RegisterReply is the server's response to a join attempt
type RegisterReply struct {
	Success bool
	Message string
}

// SendMessageArgs contains message parameters
type SendMessageArgs struct {
	Sender  string
	Content string
}

// GetMessagesArgs requests unread messages
type GetMessagesArgs struct {
	Username  string
	LastIndex int
}

// Message represents a single chat message
type Message struct {
	Sender   string
	Content  string
	IsSystem bool
}

// GetMessagesReply contains new messages for the client
type GetMessagesReply struct {
	Messages []Message
	NewIndex int
}

// Void is used for empty replies
type Void struct{}
