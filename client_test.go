package infobip_test

import (
	"fmt"
	"github.com/gaart/go-infobip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *infobip.Client
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	var err error

	client, err = infobip.New(infobip.BaseURL(server.URL))
	if err != nil {
		panic(err.Error())
	}

	// fake auth endpoint
	mux.HandleFunc("/auth/1/session", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("auth-success-response.json"))
	})

	// fake sms sending endpoint
	mux.HandleFunc("/sms/1/text/single", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("sms-sent-response.json"))
	})

	// fake delivery report endpoint
	mux.HandleFunc("/sms/1/reports", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("delivery-report-response.json"))
	})

	return func() {
		server.Close()
	}
}

func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestEmptyCredentialsClient(t *testing.T) {
	_, err := infobip.NewClient("", "")
	if err == nil {
		t.Fatal("Should fail without credentials")
	}
}

func TestInvalidCredentialsClient(t *testing.T) {
	_, err := infobip.NewClient("invalid", "invalid")
	if err == nil {
		t.Fatal("Should fail without credentials")
	}
}

func TestSmsClientOnFakeAPI(t *testing.T) {

	useRealApi := os.Getenv("INFOBIP_TEST_USE_REAL_API") == "1"
	if useRealApi {
		// run this test case on fake api only
		return
	}

	username := "fake"
	password := "fake"
	testPhoneNumber := "+12125551234"

	tearDown := setup()
	defer tearDown()

	err := client.Authenticate(username, password)
	if err != nil {
		t.Fatal(err.Error())
	}

	sms := infobip.SMS{
		From: "from somebody",
		To:   []string{testPhoneNumber},
		Text: "some message",
	}

	res, err := client.SendSMS(&sms)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(res.Messages) != 1 {
		t.Fatalf("no messages in response: %+v", res)
	}

	if len(res.Messages[0].MessageID) < 1 {
		t.Fatalf("no message ID: %+v", res)
	}

	messageID := res.Messages[0].MessageID

	deliveryResults, err := client.GetDeliveryReport(messageID)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(deliveryResults.Results) < 1 {
		t.Fatalf("Must include at least one delivery report")
	}

	if len(deliveryResults.Results[0].MessageID) < 1 {
		t.Fatalf("Message ID not found for delivery report")
	}

	if deliveryResults.Results[0].Price.PricePerMessage.String() != "1.23" {
		t.Fatalf("Price parsing failed")
	}
}

func TestSmsClientOnRealAPI(t *testing.T) {

	useRealApi := os.Getenv("INFOBIP_TEST_USE_REAL_API") == "1"
	if !useRealApi {
		// run this test case on real api only
		return
	}

	var err error
	client, err = infobip.New()
	if err != nil {
		panic(err.Error())
	}

	username := os.Getenv("INFOBIP_USERNAME")
	password := os.Getenv("INFOBIP_PASSWORD")
	testPhoneNumber := os.Getenv("INFOBIP_TEST_PHONE_NUMBER")

	err = client.Authenticate(username, password)
	if err != nil {
		t.Fatal(err.Error())
	}

	sms := infobip.SMS{
		From: "from somebody",
		To:   []string{testPhoneNumber},
		Text: "some message",
	}

	res, err := client.SendSMS(&sms)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(res.Messages) != 1 {
		t.Fatalf("no messages in response: %+v", res)
	}

	if len(res.Messages[0].MessageID) < 1 {
		t.Fatalf("no message ID: %+v", res)
	}

	time.Sleep(10 * time.Second)

	messageID := res.Messages[0].MessageID

	deliveryResults, err := client.GetDeliveryReport(messageID)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(deliveryResults.Results) < 1 {
		t.Fatalf("Must include at least one delivery report")
	}

	if len(deliveryResults.Results[0].MessageID) < 1 {
		t.Fatalf("Message ID not found for delivery report")
	}

}
