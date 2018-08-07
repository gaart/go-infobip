package infobip

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestEmptyCredentialsClient(t *testing.T) {
	_, err := NewClient("", "")
	if err == nil {
		t.Fatal("Should fail without credentials")
	}
}

func TestInvalidCredentialsClient(t *testing.T) {
	_, err := NewClient("invalid", "invalid")
	if err == nil {
		t.Fatal("Should fail without credentials")
	}
}

func TestSmsClient(t *testing.T) {

	username := os.Getenv("INFOBIP_USERNAME")
	password := os.Getenv("INFOBIP_PASSWORD")
	testPhoneNumber := os.Getenv("INFOBIP_TEST_PHONE_NUMBER")

	client, err := NewClient(username, password)
	if err != nil {
		t.Fatal(err.Error())
	}

	sms := SMS{
		"from somebody",
		[]string{testPhoneNumber},
		"some message",
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

	fmt.Printf("\n%v", res)

	deliveryResults, err := client.GetDeliveryReport(messageID)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("\n%+v", deliveryResults)

	if len(deliveryResults.Results) < 1 {
		t.Fatalf("Must include at least one delivery report")
	}

	if len(deliveryResults.Results[0].MessageID) < 1 {
		t.Fatalf("Message ID not found for delivery report")
	}
}
