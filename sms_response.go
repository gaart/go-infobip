package infobip

// SmsResponseStatus indicates whether the message is successfully sent, not sent,
// delivered, not delivered, waiting for delivery or any other possible status.
type SmsResponseStatus struct {
	GroupID     int    `json:"groupId"`
	GroupName   string `json:"groupName"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SmsResponseDetails contains info about every message.
type SmsResponseDetails struct {
	To        string            `json:"to"`
	Status    SmsResponseStatus `json:"status"`
	SmsCount  int               `json:"smsCount"`
	MessageID string            `json:"messageId"`
}

// SmsResponse contains an array of sent message objects, one object per every message.
// BulkID is the ID that uniquely identifies the request.
// Bulk ID will be received only when you send a message to more than one destination address.
type SmsResponse struct {
	BulkID   string               `json:"bulkId"`
	Messages []SmsResponseDetails `json:"messages"`
}
