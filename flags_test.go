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

func TestSetFlagValue(t *testing.T) {
	flags := FlagValues{}
	expectedValue := "richard"
	flags.SetValue("option", "person", expectedValue)
	value, _ := flags.GetValue("option", "person")
	if !reflect.DeepEqual(value, expectedValue) {
		t.Errorf("- Expected Value: %v", expectedValue)
		t.Errorf("- Actual Value: %v", value)
	}
}

func TestGetFlagsByType(t *testing.T) {
	flags := FlagValues{
		"option:user": "richard",
		"option:log":  "console",
		"setting:url": "www",
	}

	var tests = []struct {
		input    string
		expected []string
	}{
		{
			input:    "option",
			expected: []string{"user", "log"},
		},
		{
			input:    "setting",
			expected: []string{"url"},
		},
		{
			input:    "unknown",
			expected: []string{},
		},
	}

	for _, test := range tests {
		result := flags.GetFlagsByType(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input %q:\n", test.input)
			t.Errorf("- Expected: %v", test.expected)
			t.Errorf("- Actual: %v", result)
		}
	}
}
func TestGetFlag(t *testing.T) {
	flags := Flags{
		Flag{
			Key:       "option:test",
			FlagType:  "option",
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
		flagType, flagName := splitFlagKey(test.input)
		flag, found := flags.getFlag(flagType, flagName)
		if !reflect.DeepEqual(flag.Name, test.expectedFlag) || !reflect.DeepEqual(found, test.expectedFound) {
			t.Errorf("Input %q:\n", test.input)
			t.Errorf("- Expected Value: %v", test.expectedFlag)
			t.Errorf("- Actual Value: %v", flag)
			t.Errorf("- Expected Found: %v", test.expectedFound)
			t.Errorf("- Actual Found: %v", found)
		}
	}
}
