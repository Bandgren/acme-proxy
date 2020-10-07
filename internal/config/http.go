package config

import (
	"net/http"
	"time"
)

func ProvideHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}
