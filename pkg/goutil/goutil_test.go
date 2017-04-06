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

func TestInt32Addr(t *testing.T) {
	tests := []struct {
		input  int32
		output *int32
	}{
		{1, Int32Addr(1)},
		{2, Int32Addr(2)},
	}

	for _, test := range tests {
		if test.input != *test.output {
			t.Errorf("The input '%d' and output '%d' doesn't match.", test.input, test.output)
		}
	}
}
