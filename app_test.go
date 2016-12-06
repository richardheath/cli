package cli

import (
	"reflect"
	"testing"
)

func TestSplitArgs(t *testing.T) {

	var tests = []struct {
		input            []string
		expectedCommands []string
		expectedFlags    []string
	}{
		{
			input:            []string{"command", "--testArg", "argValue", "sub"},
			expectedCommands: []string{"command", "sub"},
			expectedFlags:    []string{"--testArg", "argValue"},
		},
		{
			input:            []string{"command", "--arg1", "arg1", "--arg2"},
			expectedCommands: []string{"command"},
			expectedFlags:    []string{"--arg1", "arg1", "--arg2"},
		},
	}

	for _, test := range tests {
		app := NewApp("test", "1.0")
		app.FlagPrefixes = []FlagPrefix{
			FlagPrefix{
				Key:         "--",
				Shorthand:   "-",
				Description: "options",
			},
		}

		app.ProcessArgs(test.input)
		if !reflect.DeepEqual(app.CommandArgs, test.expectedCommands) || !reflect.DeepEqual(app.FlagArgs, test.expectedFlags) {
			t.Errorf("Input %q:\n", test.input)
			t.Errorf("- Expected Commands: %v", test.expectedCommands)
			t.Errorf("- Actual Commands: %v", app.CommandArgs)
			t.Errorf("- Expected Flags: %v", test.expectedFlags)
			t.Errorf("- Actual Flags: %v", app.FlagArgs)
		}
	}
}
