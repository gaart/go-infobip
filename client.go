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

// Client is the top-level client.
type Client struct {
	authenticator Auth
	BaseURL       string
}

// NewClient is the constructor for the Client.
func NewClient(username, password string) (*Client, error) {

	if len(username) < 1 || len(password) < 1 {
		return nil, errors.New("username and password must be specified")
	}

	client := &Client{
		BaseURL: "https://api.infobip.com",
	}

	err := client.authenticate(username, password)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) doRequest(method string, path string, payload io.Reader, result interface{}) error {

	req, err := http.NewRequest(method, path, payload)
	if err != nil {
		return err
	}

	if len(c.authenticator.Token) > 0 {
		c.authenticator.setAuth(req)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("User-Agent", "go-infobip/0.1")

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

	return err
}

func (c *Client) authenticate(username, password string) error {

	data, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return err
	}

	res := Auth{}
	err = c.doRequest("POST", c.BaseURL+sessionEndpoint, bytes.NewBuffer(data), &res)
	if err != nil {
		return err
	}

	c.authenticator = res

	return nil
}

// GetDeliveryReport allows you to get one time delivery reports for sent SMS.
func (c *Client) GetDeliveryReport(smsID string) (*SmsReportResponse, error) {

	res := SmsReportResponse{}
	err := c.doRequest("GET", c.BaseURL+reportsEndpoint+"?messageId="+smsID, nil, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// SendSMS allows you to send a single textual message to array of destination addresses.
func (c *Client) SendSMS(sms *SMS) (*SmsResponse, error) {

	res := SmsResponse{}
	err := c.doRequest("POST", c.BaseURL+smsEndpoint, sms.buffer(), &res)
	if err != nil {
		return nil, err
	}

	if len(res.Messages) < 1 {
		return nil, errors.Errorf("Couldn't send a message: %+v", res)
	}

	return &res, nil
}
