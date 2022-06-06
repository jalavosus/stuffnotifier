package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func HttpClientWithTimeout(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func BuildRequestEndpoint(baseEndpoint, path string) string {
	return fmt.Sprintf(
		"https://%[1]s/%[2]s",
		strings.TrimSuffix(baseEndpoint, "/"),
		strings.TrimPrefix(path, "/"),
	)
}
