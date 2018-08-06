package infobip

type SmsResponseStatus struct {
	GroupId     int    `json:"groupId"`
	GroupName   string `json:"groupName"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SmsResponseDetails struct {
	To        string            `json:"to"`
	Status    SmsResponseStatus `json:"status"`
	SmsCount  int               `json:"smsCount"`
	MessageId string            `json:"messageId"`
}

type SmsResponse struct {
	BulkId   string               `json:"bulkId"`
	Messages []SmsResponseDetails `json:"messages"`
}
