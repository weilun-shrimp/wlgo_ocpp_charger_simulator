package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandleStatus(t *testing.T) {
	t.Run("usage OCPP 1.6", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handleStatus(ctx, nil)
		out := buf.String()
		if !strings.Contains(out, "Usage: status <status>") {
			t.Errorf("missing usage line: %q", out)
		}
		if !strings.Contains(out, "Valid statuses (OCPP 1.6):") || !strings.Contains(out, "SuspendedEV") {
			t.Errorf("missing 1.6 valid-status list: %q", out)
		}
	})

	t.Run("usage OCPP 2.0.1", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg201())
		handleStatus(ctx, nil)
		out := buf.String()
		if !strings.Contains(out, "Valid statuses (OCPP 2.0.1):") || !strings.Contains(out, "Occupied") {
			t.Errorf("missing 2.0.1 valid-status list: %q", out)
		}
	})

	t.Run("success", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handleStatus(ctx, []string{"Charging"})
		if !strings.Contains(buf.String(), "Status updated to: Charging") {
			t.Errorf("got %q", buf.String())
		}
		if f.status != "Charging" {
			t.Errorf("status not set, got %q", f.status)
		}
	})

	t.Run("error", func(t *testing.T) {
		f := &fakeCharger{setStatusErr: errors.New("invalid status")}
		ctx, buf := newCtx(f, cfg16())
		handleStatus(ctx, []string{"Bogus"})
		if !strings.Contains(buf.String(), "Error: invalid status") {
			t.Errorf("got %q", buf.String())
		}
	})
}
