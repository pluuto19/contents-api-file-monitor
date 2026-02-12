package main

import (
	"contents-api-file-monitor/internal/config"
	"contents-api-file-monitor/internal/logger"
	"contents-api-file-monitor/internal/requests"
	"contents-api-file-monitor/internal/twilio"
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const ALERT_MESSAGE string = "README file changed upstream"

// TODO: handle SIGHUP for reloading env vars live
func main() {
	log := logger.NewLogger(nil)
	logger.Info(log, "Starting file monitor")

	logger.Info(log, "Loading runtime configuration")
	vars := config.LoadRuntimeConfig(log)
	logger.Infof(log, "Configuration loaded")

	logger.Infof(log, "Creating HTTP client")
	client := requests.NewHTTPClient(log, time.Duration(vars.ClientTimeoutSec)*time.Second)
	logger.Info(log, "HTTP client created")

	logger.Infof(log, "Creating Twilio client")
	tc := twilio.NewWhatsAppClient(vars)
	logger.Infof(log, "Twilio client created")

	ctx, stopFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopFunc()

	loopTicker := time.NewTicker(time.Duration(60/vars.ReqFreq) * time.Minute)

	startMainLoop(ctx, client, tc, loopTicker, vars, log)
	logger.Info(log, "Application shutting down")
}

func startMainLoop(ctx context.Context, client *http.Client, tc *twilio.TwilioClient, t *time.Ticker, vars *config.RuntimeVars, log *logger.Logger) {
	var eTag, hash string
	logger.Info(log, "Main loop initialized")

	for {
		select {
		case <-t.C:
			logger.Infof(log, "Tick received, sending GET request")
			status, newETag, body, err := requests.SendGETRequest(client, ctx, vars.FileUrl, eTag, log)
			if err != nil {
				logger.ErrorWithErr(log, "Failed to send GET request", err)
				continue
			}

			switch status {
			case http.StatusOK:
				logger.Info(log, "File content has changed upstream (Status 200)")
				eTag = newETag

				if body != nil && hash != body.Sha {
					hash = body.Sha

					err = twilio.SendMessage(tc, ALERT_MESSAGE)
					if err != nil {
						logger.ErrorWithErr(log, "Error sending whatsapp message", err)
						continue
					}

				} else if body != nil {
					logger.Infof(log, "Hash unchanged: %s", hash)
				}

			case http.StatusNotModified:
				logger.Info(log, "File content unchanged (Status 304)")

			default:
				logger.Errorf(log, "Unexpected status code received: %d", status)
				continue
			}

		case <-ctx.Done():
			logger.Info(log, "Context cancelled, stopping main loop")
			return
		}
	}
}