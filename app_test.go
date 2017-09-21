package cli

import (
	"reflect"
	"testing"
)

func TestSplitArguments(t *testing.T) {
	var prefixes = map[string]string{
		"--": "option",
		"-":  "option",
	}

	var tests = []struct {
		input            []string
		expectedCommands []string
		expectedFlags    FlagValues
	}{
		{
			input:            []string{"command", "--testArg", "argValue", "sub"},
			expectedCommands: []string{"command", "sub"},
			expectedFlags:    FlagValues{"option:testArg": "argValue"},
		},
		{
			input:            []string{"command", "--arg1", "argValue", "--arg2"},
			expectedCommands: []string{"command"},
			expectedFlags:    FlagValues{"option:arg1": "argValue", "option:arg2": "true"},
		},
	}

	for _, test := range tests {
		commandArgs, FlagValues := splitPathAndFlagValues(prefixes, test.input)
		if !reflect.DeepEqual(commandArgs, test.expectedCommands) || !reflect.DeepEqual(FlagValues, test.expectedFlags) {
			t.Errorf("Input %q:\n", test.input)
			t.Errorf("- Expected Commands: %v", test.expectedCommands)
			t.Errorf("- Actual Commands:   %v", commandArgs)
			t.Errorf("- Expected Flags:    %v", test.expectedFlags)
			t.Errorf("- Actual Flags:      %v", FlagValues)
		}
	}
}

func TestMultiLevelExecuteChain(t *testing.T) {
	var log = []string{}
	var expectedLog = []string{"greet", "richard", "person"}

	var app = NewApp("test", "1.0.0")
	app.FlagType("option", "--", "-")
	var greet = app.Command("greet", "", func(ctx *Context) error {
		log = append(log, "greet")
		return nil
	})

	var greetPerson = greet.Command("person {{option:person}}", "", func(ctx *Context) error {
		log = append(log, "person")
		return nil
	})
	greetPerson.Flag("option", "person p", "", func(value string, ctx *Context) error {
		log = append(log, value)
		return nil
	})

	var context = Context{
		App:   &app,
		Flags: FlagValues{},
	}

	var commandPath = []string{"greet", "person", "richard"}
	var err = executeCommandChain(*(app.baseCommand), commandPath, &context)
	if err != nil {
		t.Errorf("Input %q:\n", commandPath)
		t.Errorf("%q", err)
	}

	if !reflect.DeepEqual(log, expectedLog) {
		t.Errorf("Input %q:\n", commandPath)
		t.Errorf("- Expected Commands: %v", expectedLog)
		t.Errorf("- Actual Commands: %v", log)
	}
}
