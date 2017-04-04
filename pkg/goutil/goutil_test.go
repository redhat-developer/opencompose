package goutil

import "testing"

func TestStringAddr(t *testing.T) {
	var _ *string = StringAddr("")
}

func TestStringOrEmpty(t *testing.T) {
	tests := []struct {
		ps *string
		s  string
	}{
		{nil, ""},
		{StringAddr(""), ""},
		{StringAddr("aaa"), "aaa"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if r := StringOrEmpty(tt.ps); r != tt.s {
				t.Fatalf("Expected %#v, got %#v", tt.s, r)
			}
		})
	}
}
