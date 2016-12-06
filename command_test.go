package cli

import (
	"reflect"
	"testing"
)

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
		app := NewApp("test", "0.1.0")
		app.FlagPrefixes = []FlagPrefix{
			FlagPrefix{
				Key:         "--",
				Shorthand:   "-",
				Description: "settings",
			},
		}
		app.Commands = []Command{
			command1,
			command2,
		}
		app.CommandArgs = test.input

		matchInfo := CommandMatchInfo{
			unprocessedArgs: app.CommandArgs[:],
			BindedFlags:     []string{},
			FlagTypes:       []FlagType{},
		}

		if getMatchingCommand(app.Commands, &matchInfo); !reflect.DeepEqual(matchInfo.Command, test.want) || !reflect.DeepEqual(matchInfo.BindedFlags, test.expectedFlags) {
			t.Errorf("Input %q:\n  E %v\n  G %v", test.input, test.want, matchInfo.Command)
			t.Errorf("Expected flags:\n  E %v\n  G %v", test.expectedFlags, matchInfo.BindedFlags)
		}
	}
}
