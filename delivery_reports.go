package infobip

import (
	"github.com/shopspring/decimal"
)

type SentSmsStatus struct {
	GroupId     int    `json:"groupId"`
	GroupName   string `json:"groupName"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

type SentSmsPrice struct {
	PricePerMessage decimal.Decimal `json:"pricePerMessage"`
	Currency        string          `json:"currency"`
}

type SentSmsError struct {
	GroupId     int    `json:"groupId"`
	GroupName   string `json:"groupName"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permanent   bool   `json:"permanent"`
}

type SentSmsReport struct {
	BulkId    string        `json:"bulkId"`
	To        string        `json:"to"`
	SentAt    string        `json:"sentAt"`
	DoneAt    string        `json:"doneAt"`
	Status    SentSmsStatus `json:"status"`
	SmsCount  int           `json:"smsCount"`
	MessageId string        `json:"messageId"`
	MccMnc    string        `json:"mccMnc"`
	Price     SentSmsPrice  `json:"price"`
	Error     SentSmsError  `json:"error"`
}

type SmsReportResponse struct {
	Results []SentSmsReport `json:"results"`
}
