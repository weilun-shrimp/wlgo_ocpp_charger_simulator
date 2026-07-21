package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandleUnplug(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		f := &fakeCharger{status: "Preparing", licensePlate: "ABC123"}
		ctx, buf := newCtx(f, cfg16())
		handleUnplug(ctx, nil)
		if !strings.Contains(buf.String(), "Car unplugged (Available)") {
			t.Errorf("got %q", buf.String())
		}
		if f.status != "Available" {
			t.Errorf("status not reset, got %q", f.status)
		}
	})

	t.Run("error", func(t *testing.T) {
		f := &fakeCharger{unplugErr: errors.New("unplug failed")}
		ctx, buf := newCtx(f, cfg16())
		handleUnplug(ctx, nil)
		if !strings.Contains(buf.String(), "Error: unplug failed") {
			t.Errorf("got %q", buf.String())
		}
	})
}
