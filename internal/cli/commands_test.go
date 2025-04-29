package cli

import (
	"testing"
)

func TestCommands(t *testing.T) {
	// Initialization test
	initializedCommands := NewCommands()

	if initializedCommands.AllCommands == nil {
		t.Errorf("commands functions map not initialized by NewCommands()")
	}
	if initializedCommands.Descriptions == nil {
		t.Errorf("commands descriptions map not initialized by NewCommands()")
	}

	// Registration test
	initializedCommands.Register("test", "test command", handlerLogin)
	if _, exists := initializedCommands.AllCommands["test"]; !exists {
		t.Errorf("command not found in functions map after attempted registration")
	}
	if _, exists := initializedCommands.Descriptions["test"]; !exists {
		t.Errorf("command not found in descriptions map after attempted registration")
	}

	// Get test
	descriptionsMap := initializedCommands.Get()
	if descriptionsMap == nil {
		t.Errorf("descriptions map not initialized by Get()")
	}
	if descriptionsMap["test"] != "test command" {
		t.Errorf("test command's description not found in descriptions map after Get()")
	}
}
