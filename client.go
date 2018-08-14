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

const apiURL = "https://api.infobip.com"

// Client is the top-level client.
type Client struct {
	authenticator Auth
	baseURL       string
}

// Option is a functional option for configuring the API client
type Option func(*Client) error

// BaseURL allows overriding of API client baseURL for testing
func BaseURL(baseURL string) Option {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

// parseOptions parses the supplied options functions and returns a configured
// *Client instance
func (c *Client) parseOptions(opts ...Option) error {
	// Range over each options function and apply it to our API type to
	// configure it. Options functions are applied in order, with any
	// conflicting options overriding earlier calls.
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}

	return nil
}

// New creates a new API client
func New(opts ...Option) (*Client, error) {

	client := &Client{
		baseURL: apiURL,
	}

	if err := client.parseOptions(opts...); err != nil {
		return nil, err
	}

	return client, nil
}

// NewClient is the constructor for the Client.
func NewClient(username, password string) (*Client, error) {

	if len(username) < 1 || len(password) < 1 {
		return nil, errors.New("username and password must be specified")
	}

	client, err := New()
	if err != nil {
		return nil, err
	}

	err = client.Authenticate(username, password)
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

// Authenticate allows you to get access token.
func (c *Client) Authenticate(username, password string) error {

	if len(username) < 1 || len(password) < 1 {
		return errors.New("username and password must be specified")
	}

	data, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		return err
	}

	res := Auth{}
	err = c.doRequest("POST", c.baseURL+sessionEndpoint, bytes.NewBuffer(data), &res)
	if err != nil {
		return err
	}

	c.authenticator = res

	return nil
}

// GetDeliveryReport allows you to get one time delivery reports for sent SMS.
func (c *Client) GetDeliveryReport(smsID string) (*SmsReportResponse, error) {

	res := SmsReportResponse{}
	err := c.doRequest("GET", c.baseURL+reportsEndpoint+"?messageId="+smsID, nil, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// SendSMS allows you to send a single textual message to array of destination addresses.
func (c *Client) SendSMS(sms *SMS) (*SmsResponse, error) {

	res := SmsResponse{}
	err := c.doRequest("POST", c.baseURL+smsEndpoint, sms.buffer(), &res)
	if err != nil {
		return nil, err
	}

	if len(res.Messages) < 1 {
		return nil, errors.Errorf("Couldn't send a message: %+v", res)
	}

	return &res, nil
}
