package infobip

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

const reportsEndpoint = "/sms/1/reports"
const smsEndpoint = "/sms/1/text/single"
const sessionEndpoint = "/auth/1/session"

type Client struct {
	authenticator Auth
	BaseURL       string
}

func NewClient(username, password string) (*Client, error) {

	if len(username) < 1 || len(password) < 1 {
		return nil, errors.New("username and password must be specified")
	}

	client := &Client{
		BaseURL: "https://api.infobip.com",
	}

	err := client.Authenticate(username, password)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) DoRequest(method string, path string, payload io.Reader, result interface{}) error {

	req, err := http.NewRequest(method, path, payload)
	if err != nil {
		return err
	}

	if len(c.authenticator.Token) > 0 {
		c.authenticator.SetAuth(req)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("User-Agent", "infobip-api-go-client/0.1")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		resp.Body.Close()
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return errors.New(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(result)

	return nil
}

func (c *Client) Authenticate(username, password string) error {

	data, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return err
	}

	res := Auth{}
	err = c.DoRequest("POST", c.BaseURL+sessionEndpoint, bytes.NewBuffer(data), &res)
	if err != nil {
		return err
	}

	c.authenticator = res

	return nil
}

func (c *Client) GetDeliveryReport(smsId string) (*SmsReportResponse, error) {

	res := SmsReportResponse{}
	err := c.DoRequest("GET", c.BaseURL+reportsEndpoint+"?messageId="+smsId, nil, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) SendSMS(sms *SMS) (*SmsResponse, error) {

	res := SmsResponse{}
	err := c.DoRequest("POST", c.BaseURL+smsEndpoint, sms.Buffer(), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
