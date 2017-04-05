package util

import (
	"testing"
	"time"
)

const (
	retryAttempts = 3
)

func TestGetURLData(t *testing.T) {
	tests := []struct {
		name    string
		succeed bool
		url     string
	}{
		{"passing a valid URL", true, "http://example.com"},
		{"passing an invalid URL", false, "invalid.url.^&*!@#"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetURLData(tt.url, retryAttempts)

			// error occurred but test was expected to pass
			if err != nil && tt.succeed {
				t.Fatalf("GetURLData was expected to succeed for the URL: %v, but it failed with the error: %v", tt.url, err)
			}

			// no error occurred but test was expected to fail
			if err == nil && !tt.succeed {
				t.Fatalf("GetURLData was expected to fail for the URL: %v, but it passed", tt.url)
			}
		})
	}
}

func TestFetchURLWithRetries(t *testing.T) {
	tests := []struct {
		name    string
		succeed bool
		url     string
	}{
		{"passing a valid URL", true, "https://example.com/"},
		{"passing a URL with no DNS resolution", false, "https://invalid.example/"},
		{"passing a blank string as URL", false, ""},
		{"passing a URL which gives a non 200 OK status code (gives 404)", false, "https://google.com/giveme404"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := FetchURLWithRetries(tt.url, retryAttempts, time.Second)

			// error occurred but tt was expected to pass
			if err != nil && tt.succeed {
				t.Fatalf("FetchURLWithRetries was expected to succeed for the URL: %v, but it failed with the error: %v", tt.url, err)
			}

			// no error occurred but test was expected to fail
			if err == nil && !tt.succeed {
				t.Fatalf("FetchURLWithRetries was expected to fail for the URL: %v, but it passed", tt.url)
			}
		})
	}
}
