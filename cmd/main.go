package main

import (
	"contents-api-file-monitor/internal/config"
	"contents-api-file-monitor/internal/logger"
	"contents-api-file-monitor/internal/requests"
	"contents-api-file-monitor/internal/twilio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ALERT_MESSAGE string = "README file changed upstream"

// TODO: handle SIGHUP for reloading env vars live
func main() {
	os.Exit(run())
}

func run() int {
	log := logger.NewLogger(nil)
	logger.Info(log, "Starting file monitor")

	logger.Info(log, "Loading runtime configuration")
	vars := config.LoadRuntimeConfig(log)
	if vars == nil {
		logger.Error(log, "Runtime configuration is nil")
		return 1
	}
	logger.Infof(log, "Configuration loaded")

	logger.Infof(log, "Creating HTTP client")
	client := requests.NewHTTPClient(time.Duration(vars.ClientTimeoutSec) * time.Second)
	if client == nil {
		logger.Error(log, "HTTP client is nil")
		return 1
	}
	logger.Info(log, "HTTP client created")

	logger.Infof(log, "Creating Twilio client")
	tc := twilio.NewWhatsAppClient(vars.TUsername, vars.TAuthTok, vars.TFrom, vars.TTo, vars.TContentSid)
	if tc == nil {
		logger.Error(log, "Twilio client is nil")
		return 1
	}
	logger.Infof(log, "Twilio client created")

	ctx, stopFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopFunc()

	loopTicker := time.NewTicker(time.Duration(60/vars.ReqFreq) * time.Minute)
	defer loopTicker.Stop()

	if err := startMainLoop(ctx, client, tc, loopTicker, vars, log); err != nil {
		logger.ErrorWithErr(log, "Main loop aborted", err)
		return 1
	}

	logger.Info(log, "Application shutting down")
	return 0
}

func startMainLoop(ctx context.Context, client *http.Client, tc *twilio.TwilioClient, t *time.Ticker, vars *config.RuntimeVars, log *logger.Logger) error {
	if ctx == nil {
		return fmt.Errorf("ctx is nil")
	}
	if client == nil {
		return fmt.Errorf("http client is nil")
	}
	if tc == nil {
		return fmt.Errorf("twilio client is nil")
	}
	if t == nil {
		return fmt.Errorf("ticker is nil")
	}
	if vars == nil {
		return fmt.Errorf("runtime vars is nil")
	}

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

					err = twilio.SendMessage(log, tc, ALERT_MESSAGE)
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
			client.CloseIdleConnections()
			return nil
		}
	}
}
