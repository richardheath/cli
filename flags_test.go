package cli

import (
	"reflect"
	"testing"
)

func TestProcessFlags(t *testing.T) {
	flagTypes := []FlagType{
		FlagType{
			Key:        "flag1",
			Shorthand:  "f",
			Default:    "flag1Default",
			Group:      "--",
			Validators: []FlagValidator{},
		},
	}
	flagPrefixes := map[string]string{
		"--": "--",
		"-":  "--",
	}
	var tests = []struct {
		input           []string
		expectedKnown   map[string]string
		expectedUnknown map[string]string
	}{
		{
			input:           []string{"--testArg", "arg"},
			expectedKnown:   map[string]string{},
			expectedUnknown: map[string]string{"--testArg": "arg"},
		},
		{
			input:           []string{"--flag1", "val"},
			expectedKnown:   map[string]string{"--flag1": "val"},
			expectedUnknown: map[string]string{},
		},
	}

	for _, test := range tests {
		if known, unknown, _ := processFlags(test.input, flagTypes, flagPrefixes); !reflect.DeepEqual(known, test.expectedKnown) || !reflect.DeepEqual(unknown, test.expectedUnknown) {
			t.Errorf("Input %q:\n", test.input)
			t.Errorf("- Expected Known: %v", test.expectedKnown)
			t.Errorf("- Actual Known: %v", known)
			t.Errorf("- Expected Unknown: %v", test.expectedUnknown)
			t.Errorf("- Actual Unknown: %v", unknown)
		}
	}
}

func TestGetKeyFlagType(t *testing.T) {
	flag1 := FlagType{
		Key:        "flag1",
		Shorthand:  "f",
		Default:    "flag1Default",
		Group:      "--",
		Validators: []FlagValidator{},
	}

	flagTypes := []FlagType{
		flag1,
	}
	flagPrefixes := map[string]string{
		"--": "--",
		"-":  "--",
	}
	var tests = []struct {
		inputKey      string
		inputPrefix   string
		expectedType  FlagType
		expectedFound bool
	}{
		{
			inputKey:      "flag1",
			inputPrefix:   "--",
			expectedType:  flag1,
			expectedFound: true,
		},
	}

	for _, test := range tests {
		if flagType, found := getKeyFlagType(flagTypes, flagPrefixes, test.inputKey, test.inputPrefix); !reflect.DeepEqual(found, test.expectedFound) || !reflect.DeepEqual(flagType.Key, test.expectedType.Key) {
			t.Errorf("Input %s%s:\n", test.inputPrefix, test.inputKey)
			t.Errorf("- Expected type: %v", test.expectedType)
			t.Errorf("- Actual type: %v", flagType)
			t.Errorf("- Expected found: %v", test.expectedFound)
			t.Errorf("- Actual found: %v", found)
		}
	}
}
