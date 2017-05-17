package util

import (
	"reflect"
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

func TestMergeMaps(t *testing.T) {
	no_map1 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	no_map2 := map[string]string{
		"key4": "value4",
		"key5": "value5",
		"key6": "value6",
	}
	no_map3 := map[string]string{
		"key7": "value7",
		"key8": "value8",
		"key9": "value9",
	}
	o_map1 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	o_map2 := map[string]string{
		"key4": "value4",
		"key2": "value5",
		"key1": "value6",
	}
	o_map3 := map[string]string{
		"key4": "value7",
		"key5": "value8",
		"key6": "value9",
	}
	no_map12 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
		"key6": "value6",
	}
	no_map123 := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
		"key6": "value6",
		"key7": "value7",
		"key8": "value8",
		"key9": "value9",
	}
	o_map12 := map[string]string{
		"key3": "value3",
		"key4": "value4",
		"key2": "value5",
		"key1": "value6",
	}
	o_map123 := map[string]string{
		"key3": "value3",
		"key2": "value5",
		"key1": "value6",
		"key4": "value7",
		"key5": "value8",
		"key6": "value9",
	}
	o_map231 := map[string]string{
		"key4": "value7",
		"key5": "value8",
		"key6": "value9",
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	tests := []struct {
		name   string
		input  []*map[string]string
		output *map[string]string
	}{
		{
			"passing 1 map",
			[]*map[string]string{&no_map1},
			&no_map1,
		},
		{
			"merging 2 non overlapping maps",
			[]*map[string]string{&no_map1, &no_map2},
			&no_map12,
		},
		{
			"merging 3 non overlapping maps",
			[]*map[string]string{&no_map1, &no_map2, &no_map3},
			&no_map123,
		},
		{
			"merging 2 overlapping maps",
			[]*map[string]string{&o_map1, &o_map2},
			&o_map12,
		},
		{
			"merging 3 overlapping maps",
			[]*map[string]string{&o_map1, &o_map2, &o_map3},
			&o_map123,
		},
		{
			"merging 3 overlapping maps in a different order",
			[]*map[string]string{&o_map2, &o_map3, &o_map1},
			&o_map231,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergedMap := MergeMaps(tt.input...)
			if !reflect.DeepEqual(*tt.output, *mergedMap) {
				t.Fatalf("The expected output - %v - is different than the resulting merged map - %v", *tt.output, *mergedMap)
			}
		})
	}
}
