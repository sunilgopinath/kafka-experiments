package models

// OrderEvent represents an order going through stages
type OrderEvent struct {
	OrderID   string `json:"orderId"`
	UserID    string `json:"userId"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}
