package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// This function validates the URL, then tries to fetch the data with retries
// and then reads and returns the data as []byte
// Returns an error if the URL is invalid, or fetching the data failed or
// if reading the response body fails.
func GetURLData(urlString string, attempts int) ([]byte, error) {
	// Validate URL
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", urlString)
	}

	// Fetch the URL and store the response body
	data, err := FetchURLWithRetries(urlString, attempts, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed fetching data from the URL %v: %s", urlString, err)
	}

	return data, nil
}

// Try to fetch the given url string, and make <attempts> attempts at it.
// Wait for <duration> time between each try.
// This returns the data from  the response body []byte upon successful fetch
// The passed URL is not validated, so validate the URL before passing to this
// function
func FetchURLWithRetries(url string, attempts int, duration time.Duration) ([]byte, error) {
	var data []byte
	var err error

	for i := 0; i < attempts; i++ {
		var response *http.Response

		// sleep for <duration> seconds before trying again
		if i > 0 {
			time.Sleep(duration)
		}

		// retry if http.Get fails
		// if all the retries fail, then return statement at the end of the
		// function will return this err received from http.Get
		response, err = http.Get(url)
		if err != nil {
			continue
		}
		defer response.Body.Close()

		// if the status code is not 200 OK, return an error
		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unable to fetch %v, server returned status code %v", url, response.StatusCode)
		}

		// Read from the response body, ioutil.ReadAll will return []byte
		data, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("reading from response body failed: %s", err)
		}
		break
	}

	return data, err
}

// MergeMaps will merge the given maps, but it does not check for conflicts.
// In case of conflicting keys, the map that is provided later overrides the previous one.
// TODO: add to docs about use with caution bits
func MergeMaps(maps ...*map[string]string) *map[string]string {
	mergedMap := make(map[string]string)
	for _, m := range maps {
		for k, v := range *m {
			mergedMap[k] = v
		}
	}
	return &mergedMap
}
