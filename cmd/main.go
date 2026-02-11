package main

import (
	"contents-api-file-monitor/internal/config"
	"context"
	"os/signal"
	"syscall"
	"time"
)

// TODO: handle SIGHUP for reloading env vars live
func main() {
	vars := config.GetRuntimeConfig()

	ctx, stopFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopFunc()

	loopTicker := time.NewTicker(time.Duration(vars.ReqFreq/60) * time.Minute)

	go startMainLoop(loopTicker, ctx, vars)
}

func startMainLoop(t *time.Ticker, ctx context.Context, vars *config.RuntimeVars) {
	for {
		select {
			// call a request function with vars.FileUrl.
		case <-t.C:
		case <-ctx.Done():
			return
		}
	}
}