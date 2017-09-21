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
