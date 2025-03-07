package consumer

// Notification represents a user notification event
type Notification struct {
	UserID    string `json:"userId"`
	Type      string `json:"type"`     // email, sms, push
	Priority  string `json:"priority"` // high, low
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
