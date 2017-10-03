package cli

import (
	"reflect"
	"testing"
)

func TestGetFlagValue(t *testing.T) {
	flags := FlagValues{
		"option:user": "richard",
	}
	var tests = []struct {
		input         string
		expectedValue string
		expectedFound bool
	}{
		{
			input:         "user",
			expectedValue: "richard",
			expectedFound: true,
		},
		{
			input:         "unknown",
			expectedValue: "",
			expectedFound: false,
		},
	}

	for _, test := range tests {
		value, found := flags.GetValue("option", test.input)
		if !reflect.DeepEqual(value, test.expectedValue) || !reflect.DeepEqual(found, test.expectedFound) {
			t.Errorf("Input %q:\n", test.input)
			t.Errorf("- Expected Value: %v", test.expectedValue)
			t.Errorf("- Actual Value: %v", value)
			t.Errorf("- Expected Found: %v", test.expectedFound)
			t.Errorf("- Actual Found: %v", found)
		}
	}
}

func TestGetFlag(t *testing.T) {
	flags := Flags{
		Flag{
			Key:       "option:test",
			Name:      "test",
			Shorthand: "t",
		},
	}
	var tests = []struct {
		input         string
		expectedFlag  string
		expectedFound bool
	}{
		{
			input:         "option:test",
			expectedFlag:  "test",
			expectedFound: true,
		},
		{
			input:         "option:unknown",
			expectedFlag:  "",
			expectedFound: false,
		},
	}

	for _, test := range tests {
		flag, found := flags.getFlag(test.input)
		if !reflect.DeepEqual(flag.Name, test.expectedFlag) || !reflect.DeepEqual(found, test.expectedFound) {
			t.Errorf("Input %q:\n", test.input)
			t.Errorf("- Expected Value: %v", test.expectedFlag)
			t.Errorf("- Actual Value: %v", flag)
			t.Errorf("- Expected Found: %v", test.expectedFound)
			t.Errorf("- Actual Found: %v", found)
		}
	}
}
