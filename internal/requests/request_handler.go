package requests

import (
	"contents-api-file-monitor/internal/dtos"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func SendGETRequest(c *http.Client, ctx context.Context, url, currETag string) (int, string, *dtos.ReadmeResponseDTO, error) {
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
		return -1, "", nil, err
	}

	if currETag != "" {
		req.Header.Set("If-None-Match", currETag)
	}

	resp, err := c.Do(req)
	if err != nil {
		return -1, "", nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	eTag := resp.Header.Get("ETag")

	if statusCode == http.StatusNotModified {
		return statusCode, eTag, nil, nil
	}

	var res dtos.ReadmeResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		if errors.Is(err, io.EOF) {
			return -1, "", nil, fmt.Errorf("unexpected empty body for status %d", statusCode)
		}
		return -1, "", nil, fmt.Errorf("decoding response: %w", err)
	}

	return statusCode, eTag, &res, nil
}