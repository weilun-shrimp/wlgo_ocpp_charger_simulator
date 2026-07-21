package cli

import (
	"strings"
	"testing"
)

func TestHandleHelp(t *testing.T) {
	t.Run("OCPP 1.6", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handleHelp(ctx, nil)
		out := buf.String()
		for _, want := range []string{
			"Available commands:",
			"connect",
			"disconnect",
			"current <amps>    - Set charging current (0-32.0 A, 0 = SuspendedEVSE)",
			"Valid statuses (OCPP 1.6):",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q in help output: %q", want, out)
			}
		}
	})

	t.Run("OCPP 2.0.1", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg201())
		handleHelp(ctx, nil)
		out := buf.String()
		if !strings.Contains(out, "Valid statuses (OCPP 2.0.1):") || !strings.Contains(out, "Occupied") {
			t.Errorf("missing 2.0.1 status list: %q", out)
		}
		if strings.Contains(out, "Valid statuses (OCPP 1.6):") {
			t.Errorf("must not show 1.6 status list under 2.0.1: %q", out)
		}
	})
}
