package cli

import (
	"strings"
	"testing"
)

func TestHandleQuit(t *testing.T) {
	// Both "quit" and "exit" must be registered to the same handler and print
	// the Ctrl+C guidance.
	for _, alias := range []string{"quit", "exit"} {
		t.Run(alias, func(t *testing.T) {
			f := &fakeCharger{}
			ctx, buf := newCtx(f, cfg16())
			Dispatch(ctx, alias, nil)
			if !strings.Contains(buf.String(), "Use Ctrl+C to exit") {
				t.Errorf("alias %q: got %q", alias, buf.String())
			}
		})
	}
}
