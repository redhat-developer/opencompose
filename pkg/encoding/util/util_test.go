package util

import (
	"errors"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

type S struct {
	A  int  `yaml:"omitempty"`
	B  *int `yaml:"omitempty"`
	Ra int
	Rb *int
}

var integer int

func TestIsRequiredField(t *testing.T) {
	tests := []struct {
		Required bool
		Name     string
	}{
		{false, "A"},
		{false, "B"},
		{true, "Ra"},
		{true, "Rb"},
	}

	for _, tt := range tests {
		field, found := reflect.TypeOf((*S)(nil)).Elem().FieldByName(tt.Name)
		if !found {
			panic(errors.New("field does not exist"))
		}

		if IsRequiredField(field) {
			if !tt.Required {
				t.Errorf("Field %q should NOT be required!", tt.Name)
			}
		} else {
			if tt.Required {
				t.Errorf("Field %q should be required!", tt.Name)
			}
		}
	}
}

func TestIsUnset(t *testing.T) {
	tests := []struct {
		Unset bool
		I     interface{}
	}{
		{true, 0},
		{false, 42},
		{true, (*int)(nil)},
		{false, new(int)},
		{true, []int{}},
		{true, make([]int, 0, 42)},
		{false, make([]int, 42)},
	}

	for _, tt := range tests {
		if IsUnset(reflect.ValueOf(tt.I)) {
			if !tt.Unset {
				t.Errorf("IsUnset returned TRUE for %#v", tt.I)
			}
		} else {
			if tt.Unset {
				t.Errorf("IsUnset returned FALSE for %#v", tt.I)
			}
		}
	}
}

func TestValidateRequiredFields(t *testing.T) {
	tests := []struct {
		Succeed bool
		Value   interface{}
	}{
		{false, S{}},
		{false, S{
			Rb: &integer,
		}},
		{true, S{
			Ra: 42,
			Rb: &integer,
		}},
		{true, []S{}},
		{false, []S{
			{},
		}},
		{false, []S{
			{
				Rb: &integer,
			},
		}},
		{true, []S{
			{
				Ra: 42,
				Rb: &integer,
			},
		}},
		{false, struct {
			A []S `yaml:"omitempty"`
		}{
			A: []S{
				{
					Rb: &integer,
				},
			},
		}},
		{true, struct {
			A []S `yaml:"omitempty"`
		}{
			A: []S{
				{
					Ra: 42,
					Rb: &integer,
				},
			},
		}},
	}

	for _, tt := range tests {
		err := ValidateRequiredFields(tt.Value)
		if err != nil {
			if tt.Succeed {
				t.Errorf("Failed to validate %#v: %s", tt.Value, err)
			}
			continue
		}

		if !tt.Succeed {
			t.Errorf("Expected %#v to fail!", tt.Value)
			continue
		}
	}
}

func TestInterfaceToJSON(t *testing.T) {
	tests := []struct {
		input  string
		output map[string]interface{}
	}{
		{
			input:  `foo: bar`,
			output: map[string]interface{}{"foo": "bar"},
		},
		{
			input: `
one:
- 1
- 2
- 3
two: four`,
			output: map[string]interface{}{
				"two": "four",
				"one": []interface{}{1, 2, 3},
			},
		},
		{
			input: `
one:
  two:
    three: four
five: six
seven:
- eight
- nine`,
			output: map[string]interface{}{
				"one": map[string]interface{}{
					"two": map[string]interface{}{
						"three": "four",
					},
				},
				"five":  "six",
				"seven": []interface{}{"eight", "nine"},
			},
		},
	}

	for _, test := range tests {
		var input interface{}
		err := yaml.Unmarshal([]byte(test.input), &input)
		if err != nil {
			t.Fatalf("failed unmarshalling: %v", err)
		}
		gotOutput := InterfaceToJSON(input)

		if !reflect.DeepEqual(test.output, gotOutput) {
			t.Errorf("Expected: %#v\nGot: %#v", test.output, gotOutput)
		}
	}
}
