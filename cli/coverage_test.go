package cli

import "testing"

// TestRegistryCoversAllCommands asserts that every canonical command verb has a
// registered handler and that no unexpected command has been registered. Adding
// a new command without registering it (or registering an unlisted one) fails
// this test, keeping "every command is covered" enforceable.
func TestRegistryCoversAllCommands(t *testing.T) {
	want := []string{
		"help", "connect", "disconnect", "plugin", "unplug",
		"status", "start", "stop", "meter", "plate",
		"soc", "current", "power", "info", "quit", "exit",
	}

	for _, name := range want {
		if _, ok := registry[name]; !ok {
			t.Errorf("no handler registered for command %q", name)
		}
	}

	if len(registry) != len(want) {
		t.Errorf("registry has %d commands, expected %d — an unlisted command was registered or one is missing",
			len(registry), len(want))
	}
}
