package twilio

import (
	"contents-api-file-monitor/internal/logger"
	"fmt"

	t "github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioClient struct {
	client     *t.RestClient
	from, to   string
	contentSid string
}

func NewWhatsAppClient(tUsername, tAuthTok, tFrom, tTo, tContentSid string) *TwilioClient {
	if tUsername == "" || tAuthTok == "" || tFrom == "" || tTo == "" || tContentSid == "" {
		return nil
	}

	return &TwilioClient{
		client: t.NewRestClientWithParams(t.ClientParams{
			Username: tUsername,
			Password: tAuthTok,
		}),
		from:       tFrom,
		to:         tTo,
		contentSid: tContentSid,
	}
}

func SendMessage(log *logger.Logger, tc *TwilioClient, msg string) error {
	if tc == nil {
		return fmt.Errorf("twilio client is nil")
	}
	if msg == "" {
		return fmt.Errorf("message is an empty string")
	}

	p := &api.CreateMessageParams{}
	p.SetTo(tc.to)
	p.SetFrom(tc.from)
	p.SetContentSid(tc.contentSid)
	p.SetContentVariables(`{"1":"` + msg + `"}`)

	_, err := tc.client.Api.CreateMessage(p)
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	return nil
}
