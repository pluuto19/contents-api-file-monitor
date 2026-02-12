package twilio

import (
	"contents-api-file-monitor/internal/config"
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

func SendMessage(tc *TwilioClient, msg string) error {
	if tc == nil {
		return fmt.Errorf("whatsapp client is nil")
	}
	if msg == "" {
		return fmt.Errorf("message is an empty string")
	}

	p := &api.CreateMessageParams{}
	p.SetTo(tc.to)
	p.SetFrom(tc.from)
	p.SetBody(msg)

	_, err := tc.client.Api.CreateMessage(p)
	if err != nil {
		return err
	}

	return nil
}