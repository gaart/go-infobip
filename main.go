package infobip

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	SmsApiUrl = "https://api.infobip.com/sms/1/text/single"
)

type SMS struct {
	From string `json:"from"`
	To   []string `json:"to"`
	Text string `json:"text"`
}

func (s *SMS) String() string {
	text, _ := json.Marshal(s)
	return string(text)
}

type SmsApiResponse struct {
	Messages []struct {
		BulkId string `json:"bulkId"`
		To     string `json:"to"`
		Status struct {
			GroupId     int    `json:"groupId"`
			GroupName   string `json:"groupName"`
			Id          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"status"`
		SmsCount  int    `json:"smsCount"`
		MessageId string `json:"messageId"`
	} `json:"messages"`
}

func (s *SMS) Send(logger *zap.Logger) (*SmsApiResponse, error) {

	req, err := http.NewRequest("POST", SmsApiUrl, strings.NewReader(s.String()))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(os.Getenv("INFOBIP_USER"), os.Getenv("INFOBIP_PASSWORD"))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := SmsApiResponse{}
	json.Unmarshal([]byte(responseData), &res)

	logger.Info("processed SMS",
		zap.Reflect("payload", s),
		zap.Reflect("api_response", res),
	)

	return &res, nil
}
