package domain

// NotificationIngestPayload payload from Android Notification Listener
type NotificationIngestPayload struct {
	PackageName string `json:"package_name" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

// ParsedTransaction output from Regex Engine
type ParsedTransaction struct {
	Amount      float64
	Type        TransactionType
	Description string
	Merchant    string
	IsValid     bool
}
