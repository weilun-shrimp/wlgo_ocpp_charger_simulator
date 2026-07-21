package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandleCurrent(t *testing.T) {
	t.Run("usage reports present value", func(t *testing.T) {
		f := &fakeCharger{current: 16}
		ctx, buf := newCtx(f, cfg16())
		handleCurrent(ctx, nil)
		out := buf.String()
		if !strings.Contains(out, "Usage: current <amperes> (0-32.0 A, 0 = SuspendedEVSE)") {
			t.Errorf("missing usage line: %q", out)
		}
		if !strings.Contains(out, "Current: 16.0 A") {
			t.Errorf("missing present value: %q", out)
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		f := &fakeCharger{current: 16}
		ctx, buf := newCtx(f, cfg16())
		handleCurrent(ctx, []string{"abc"})
		if !strings.Contains(buf.String(), "Error: invalid current value: abc") {
			t.Errorf("got %q", buf.String())
		}
		if f.current != 16 {
			t.Errorf("current must be unchanged on invalid input, got %v", f.current)
		}
	})

	t.Run("success", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handleCurrent(ctx, []string{"10"})
		if !strings.Contains(buf.String(), "Current set to: 10.0 A") {
			t.Errorf("got %q", buf.String())
		}
		if f.current != 10 {
			t.Errorf("current not set, got %v", f.current)
		}
	})

	t.Run("error", func(t *testing.T) {
		f := &fakeCharger{setCurrentErr: errors.New("exceeds maximum")}
		ctx, buf := newCtx(f, cfg16())
		handleCurrent(ctx, []string{"999"})
		if !strings.Contains(buf.String(), "Error: exceeds maximum") {
			t.Errorf("got %q", buf.String())
		}
	})
}
