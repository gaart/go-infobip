package infobip

import (
	"encoding/json"
	"go.uber.org/zap"
	"testing"
)

func TestSmsSending(t *testing.T) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sms := SMS{}
	json.Unmarshal([]byte("{\"from\": \"somebody\",\n\"to\": \"12125551234\",\n\"text\": \"tstmsg\"}"), &sms)

	res, err := sms.Send(logger)
	if err != nil {
		t.Fail()
	}

	if len(res.Messages) != 1 {
		t.Fail()
	}

	if len(res.Messages[0].MessageId) < 1 {
		t.Fail()
	}
}
