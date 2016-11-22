package cli

import (
	"reflect"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	flagPrefixes := map[string]string{
		"--": "--",
	}

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
		if commands, flags := splitCommandsAndFlags(test.input, flagPrefixes); !reflect.DeepEqual(commands, test.expectedCommands) || !reflect.DeepEqual(flags, test.expectedFlags) {
			t.Errorf("Input %q:\n", test.input)
			t.Errorf("- Expected Commands: %v", test.expectedCommands)
			t.Errorf("- Actual Commands: %v", commands)
			t.Errorf("- Expected Flags: %v", test.expectedFlags)
			t.Errorf("- Actual Flags: %v", flags)
		}
	}
}

func TestBasicCommandMatching(t *testing.T) {
	commandSinglePath := Command{
		Path:      []string{"single"},
		Commands:  []Command{},
		FlagTypes: nil,
		Action:    nil,
	}

	commandWithMultiPath := Command{
		Path:      []string{"path1", "path2"},
		Commands:  []Command{},
		FlagTypes: nil,
		Action:    nil,
	}

	commandWithFlagBinding := Command{
		Path:      []string{"bind", "{{--flagKey}}"},
		Commands:  []Command{},
		FlagTypes: nil,
		Action:    nil,
	}

	command1 := Command{
		Path: []string{"cmd1"},
		Commands: []Command{
			commandSinglePath,
			commandWithMultiPath,
			commandWithFlagBinding,
		},
		FlagTypes: nil,
		Action:    nil,
	}

	command2 := Command{
		Path: []string{"cmd2"},
		Commands: []Command{
			commandSinglePath,
			commandWithMultiPath,
			commandWithFlagBinding,
		},
		FlagTypes: nil,
		Action:    nil,
	}

	app := App{
		Name:    "test",
		Version: "0.1.0",
		FlagPrefixes: []FlagPrefix{
			FlagPrefix{
				Key:         "--",
				Shorthand:   "-",
				Description: "settings",
			},
		},
		FlagTypes: []FlagType{},
		Commands: []Command{
			command1,
			command2,
		},
	}

	var tests = []struct {
		input         []string
		want          Command
		expectedFlags []string
	}{
		{[]string{"cmd1"}, command1, []string{}},
		{[]string{"cmd2"}, command2, []string{}},
		{[]string{"cmd1", "single"}, commandSinglePath, []string{}},
		{[]string{"cmd1", "path1", "path2"}, commandWithMultiPath, []string{}},
		{[]string{"cmd1", "bind", "someInput"}, commandWithFlagBinding, []string{"--flagKey", "someInput"}},
		{[]string{"cmd2", "single"}, commandSinglePath, []string{}},
	}

	for _, test := range tests {
		flagArgs := make([]string, 0, 2)
		if got, _, bindedFlags, _ := getMatchingCommand(app.Commands, test.input, app.FlagTypes, flagArgs); !reflect.DeepEqual(got.Path, test.want.Path) || !reflect.DeepEqual(bindedFlags, test.expectedFlags) {
			t.Errorf("Input %q:\n  E %v\n  G %v", test.input, test.want, got)
			t.Errorf("Expected flags:\n  E %v\n  G %v", test.expectedFlags, bindedFlags)
		}
	}
}
