package main

import (
	"contents-api-file-monitor/internal/config"
	"contents-api-file-monitor/internal/requests"
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// TODO: handle SIGHUP for reloading env vars live
func main() {
	vars := config.LoadRuntimeConfig()
	client := requests.NewHTTPClient(time.Duration(vars.ClientTimeoutSec))

	ctx, stopFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopFunc()

	loopTicker := time.NewTicker(time.Duration(vars.ReqFreq/60) * time.Minute)

	go startMainLoop(client, loopTicker, ctx, vars)
}

func startMainLoop(client *http.Client, t *time.Ticker, ctx context.Context, vars *config.RuntimeVars) {
	var eTag, hash string

	for {
		select {
		case <-t.C:
			status, newETag, body, err := requests.SendGETRequest(client, ctx, vars.FileUrl, eTag)
			if err != nil {
				// TODO: Log the error and continue
			}
			if status != http.StatusNotModified && status != http.StatusOK {
				// TODO: Log the error and continue
			}

			if status == http.StatusOK {
				// Something changed upstream
				eTag = newETag
				if hash != body.Sha {
					hash = body.Sha
					// TODO: Alert using Twilio
				}
			}

		case <-ctx.Done():
			return
		}
	}
}