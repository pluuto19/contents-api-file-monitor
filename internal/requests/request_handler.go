package requests

import (
	"contents-api-file-monitor/internal/dtos"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var c *http.Client

func SetupHTTPClient(timeout time.Duration) {
	if c == nil {
		c = &http.Client{
			Timeout: timeout,
		}
	}
}

func SendGETRequest(ctx context.Context, url, currETag string) (int, string, *dtos.ReadmeResponseDTO, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return -1, "",  nil, err
	}

	if currETag != "" {
		req.Header.Set("If-None-Match", currETag)
	}

	resp, err := c.Do(req)
	if err != nil {
		return -1, "",  nil, err
	}
	defer resp.Body.Close()

	var res dtos.ReadmeResponseDTO
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&res)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return -1, "", nil, fmt.Errorf("empty body: %w", err)
		}
		return -1, "", nil, fmt.Errorf("decoding response: %w", err)
	}

	eTag := resp.Header.Get("ETag")
	status, err := strconv.ParseInt(strings.Split(resp.Status, " ")[0], 10, 32)
	if err != nil {
		return -1, "",  nil, err
	}

	return int(status), eTag, &res, nil
}