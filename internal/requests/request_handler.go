package requests

import (
	"contents-api-file-monitor/internal/dtos"
	"contents-api-file-monitor/internal/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewHTTPClient(timeout time.Duration) *http.Client {
	client := &http.Client{
		Timeout: timeout,
	}
	return client
}

func SendGETRequest(c *http.Client, ctx context.Context, url, currETag string, log *logger.Logger) (int, string, *dtos.ReadmeResponseDTO, error) {
	if c == nil {
		return -1, "", nil, fmt.Errorf("client is nil")
	}
	if url == "" {
		return -1, "", nil, fmt.Errorf("url is an empty string")
	}
	if ctx == nil {
		return -1, "", nil, fmt.Errorf("context is nil")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return -1, "", nil, fmt.Errorf("create request failed: %w", err)
	}

	if currETag != "" {
		req.Header.Set("If-None-Match", currETag)
	}

	logger.Info(log, "Sending HTTP request")
	resp, err := c.Do(req)
	if err != nil {
		return -1, "", nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	eTag := resp.Header.Get("ETag")
	logger.Infof(log, "Received response: Status: %d, ETag: %s", statusCode, eTag)

	if statusCode == http.StatusNotModified {
		logger.Info(log, "Content not modified (304), returning early")
		return statusCode, eTag, nil, nil
	}

	var res dtos.ReadmeResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		if errors.Is(err, io.EOF) {
			return -1, "", nil, fmt.Errorf("unexpected empty body for status %d", statusCode)
		}

		return -1, "", nil, fmt.Errorf("decode response failed: %w", err)
	}

	return statusCode, eTag, &res, nil
}
