package twilio

import (
	"contents-api-file-monitor/internal/config"
	"contents-api-file-monitor/internal/logger"
	"fmt"

	t "github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioClient struct {
	client *t.RestClient
	from   string
	to     string
}

func NewWhatsAppClient(vars *config.RuntimeVars) *TwilioClient {
	return &TwilioClient{
		client: t.NewRestClientWithParams(t.ClientParams{
			Username: vars.TUsername,
			Password: vars.TAuthTok,
		}),
		from: vars.TFrom,
		to:   vars.TTo,
	}
}

func SendMessage(log *logger.Logger, tc *TwilioClient, msg string) error {
	if tc == nil {
		logger.Error(log, "Twilio client is nil")
		return fmt.Errorf("twilio client is nil")
	}
	if msg == "" {
		logger.Error(log, "message is an empty string")
		return fmt.Errorf("message is an empty string")
	}

	p := &api.CreateMessageParams{}
	p.SetTo(tc.to)
	p.SetFrom(tc.from)
	p.SetBody(msg)

	_, err := tc.client.Api.CreateMessage(p)
	if err != nil {
		logger.ErrorWithErr(log, "Error while sending message", err)
		return err
	}

	return nil
}