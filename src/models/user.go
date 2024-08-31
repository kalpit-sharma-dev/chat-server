package models

// Define a struct for user status
type UserStatus struct {
	UserID string
	Online bool
}

// Define a struct for push notification payload
type NotificationPayload struct {
	To   string `json:"to"`
	Data struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"data"`
}
