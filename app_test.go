package cli

import (
	"fmt"
	"reflect"
	"strings"
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

func TestAppRun(t *testing.T) {
	var tests = []struct {
		input         []string
		description   string
		expected      []string
		expectedError string
	}{
		{
			input:       []string{"greet", "bind", "someone"},
			description: "someone value must bind to person flag",
			expected:    []string{"log", "console", "person", "someone", "greet", "bind"},
		},
		{
			input:       []string{"greet", "flag", "--person", "fullName"},
			description: "flag must use given value",
			expected:    []string{"log", "console", "person", "fullName", "greet", "flag"},
		},
		{
			input:       []string{"greet", "flag", "-p", "shortHand"},
			description: "flag must use given value using shorthand",
			expected:    []string{"log", "console", "person", "shortHand", "greet", "flag"},
		},
		{
			input:         []string{"greet", "bad_flag", "--bad", "bad"},
			description:   "flag that errors out must return error",
			expected:      []string{},
			expectedError: "bad flag",
		},
		{
			input:       []string{"greet", "default"},
			description: "flag must use default value if specified",
			expected:    []string{"log", "console", "greet", "message", "defaultHello", "default"},
		},
		{
			input:       []string{"greet", "default", "--message", "hello"},
			description: "flag must always use given value",
			expected:    []string{"log", "console", "greet", "message", "hello", "default"},
		},
		{
			input:         []string{"unknown", "command"},
			description:   "app must return not found when there's no matching command",
			expected:      []string{},
			expectedError: "Command not found",
		},
		{
			input:         []string{"greet", "bad", "command"},
			description:   "command that errors out must return error",
			expected:      []string{},
			expectedError: "bad command",
		},
		{
			input:       []string{"greet", "flag", "--person", "someone", "--log", "file"},
			description: "--log must use provided value",
			expected:    []string{"log", "file", "person", "someone", "greet", "flag"},
		},
		{
			input:       []string{"greet", "flag", "--log", "--person", "someone"},
			description: "--log should still be default value since value is not defined",
			expected:    []string{"log", "console", "person", "someone", "greet", "flag"},
		},
	}

	for _, test := range tests {
		app, log := CreateDummyApp()
		err := app.Run(test.input)
		if test.expectedError == "" && err != nil {
			t.Errorf("It should not fail: %v", err)
		}

		if test.expectedError != "" && (err == nil || !strings.Contains(err.Error(), test.expectedError)) {
			t.Errorf("Input %q:\n", test.input)
			t.Errorf("- Expected error: %v", test.expectedError)
			t.Errorf("- Actual error:   %v", err)
		}

		if test.expectedError == "" && !reflect.DeepEqual(*log, test.expected) {
			t.Errorf("%s:\n", test.description)
			t.Errorf("- Input %q:\n", test.input)
			t.Errorf("- Expected: %v", test.expected)
			t.Errorf("- Actual:   %v", log)
		}
	}
}

func TestDynamicFlags(t *testing.T) {
	app := NewApp("test", "1.0.0")
	app.FlagType("setting", "#")

	app.Command("test", func(ctx *Context) error {
		expectedFlags := []string{"dynamic"}
		expectedValue := "flag"

		flags := ctx.Flags.GetFlagsByType("setting")
		if !reflect.DeepEqual(flags, expectedFlags) {
			t.Errorf("expecting flags: %v, actual: %v", expectedFlags, flags)
		}

		value, _ := ctx.Flags.GetValue("setting", "dynamic")
		if expectedValue != value {
			t.Errorf("expecting values: %v, actual: %v", expectedValue, value)
		}
		return nil
	})

	err := app.Run([]string{"test", "#dynamic", "flag"})
	if err != nil {
		t.Error("it should not fail")
	}

}

func CreateDummyApp() (*App, *[]string) {
	var log = []string{}
	var app = NewApp("test", "1.0.0")
	app.FlagType("option", "--", "-")
	app.Flag("option", "log", "console", echoFlagAction(&log, "log", nil))

	var greet = app.Command("greet", echoCommandAction(&log, "greet", nil))
	greet.Flag("option", "person p", "", echoFlagAction(&log, "person", nil))

	greet.Command("bind {{option:person}}", echoCommandAction(&log, "bind", nil))
	greet.Command("flag", echoCommandAction(&log, "flag", nil))
	greet.Command("bad command", echoCommandAction(&log, "bind", fmt.Errorf("bad command")))

	badFlag := greet.Command("bad_flag", echoCommandAction(&log, "bad_flag", nil))
	badFlag.Flag("option", "bad", "", echoFlagAction(&log, "bad", fmt.Errorf("bad flag")))

	defaultValue := greet.Command("default", echoCommandAction(&log, "default", nil))
	defaultValue.Flag("option", "message", "defaultHello", echoFlagAction(&log, "message", nil))

	return &app, &log
}

func echoCommandAction(log *[]string, message string, err error) CommandAction {
	return func(ctx *Context) error {
		if err != nil {
			return err
		}

		*log = append(*log, message)
		return nil
	}
}

func echoFlagAction(log *[]string, message string, err error) FlagAction {
	return func(value string, ctx *Context) error {
		if err != nil {
			return err
		}

		*log = append(*log, message)
		*log = append(*log, value)
		return nil
	}
}
