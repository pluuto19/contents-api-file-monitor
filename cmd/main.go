package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stopFunc := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopFunc()

	loopTicker := time.NewTicker(15 * time.Minute)

	go startMainLoop(loopTicker, ctx)
}

func startMainLoop(t *time.Ticker, ctx context.Context) {
	for {
		select {
		case <-t.C:
		case <-ctx.Done():
			return
		}
	}
}