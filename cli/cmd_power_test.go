package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandlePower(t *testing.T) {
	t.Run("usage reports present value", func(t *testing.T) {
		f := &fakeCharger{power: 2300}
		ctx, buf := newCtx(f, cfg16())
		handlePower(ctx, nil)
		out := buf.String()
		if !strings.Contains(out, "Usage: power <watts> (0-7360.0 W, 0 = SuspendedEVSE)") {
			t.Errorf("missing usage line: %q", out)
		}
		if !strings.Contains(out, "Power: 2300.0 W") {
			t.Errorf("missing present value: %q", out)
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		f := &fakeCharger{power: 2300}
		ctx, buf := newCtx(f, cfg16())
		handlePower(ctx, []string{"abc"})
		if !strings.Contains(buf.String(), "Error: invalid power value: abc") {
			t.Errorf("got %q", buf.String())
		}
		if f.power != 2300 {
			t.Errorf("power must be unchanged on invalid input, got %v", f.power)
		}
	})

	t.Run("success", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handlePower(ctx, []string{"1000"})
		if !strings.Contains(buf.String(), "Power set to: 1000.0 W") {
			t.Errorf("got %q", buf.String())
		}
		if f.power != 1000 {
			t.Errorf("power not set, got %v", f.power)
		}
	})

	t.Run("error", func(t *testing.T) {
		f := &fakeCharger{setPowerErr: errors.New("exceeds maximum")}
		ctx, buf := newCtx(f, cfg16())
		handlePower(ctx, []string{"99999"})
		if !strings.Contains(buf.String(), "Error: exceeds maximum") {
			t.Errorf("got %q", buf.String())
		}
	})
}
