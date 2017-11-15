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
	To   string `json:"to"`
	Text string `json:"text"`

	payload string
}

type SmsApiResponse struct {
	Messages []struct {
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

func (s *SMS) setPayload() error {

	body := map[string]string{
		"from": s.From,
		"to":   s.To,
		"text": s.Text,
	}

	data, err := json.Marshal(body)
	if err == nil {
		s.payload = string(data)
	}
	return err
}

func (s *SMS) Send(logger *zap.Logger) (*SmsApiResponse, error) {

	if err := s.setPayload(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", SmsApiUrl, strings.NewReader(s.payload))
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
		zap.Reflect("api_response", res),
	)

	return &res, nil
}
