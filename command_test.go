package cli

import (
	"reflect"
	"testing"
)

func TestAddSubCommand(t *testing.T) {
	command := Command{
		Path:     []string{"test"},
		Commands: []*Command{},
	}

	subCommand := command.Command("sub cmd", "desc", nil)

	if subCommand != command.Commands[0] {
		t.Errorf("returned command doesnt match array pointer")
	}

	expectedPath := []string{"sub", "cmd"}
	if !reflect.DeepEqual(subCommand.Path, expectedPath) {
		t.Errorf("Path mismatch")
		t.Errorf("- Expected: %v", expectedPath)
		t.Errorf("- Actual: %v", subCommand.Path)
	}
}
