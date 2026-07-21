package cli

import (
	"strings"
	"testing"
)

func TestHandleDisconnect(t *testing.T) {
	t.Run("not connected", func(t *testing.T) {
		f := &fakeCharger{connected: false}
		ctx, buf := newCtx(f, cfg16())
		handleDisconnect(ctx, nil)
		if !strings.Contains(buf.String(), "Not connected") {
			t.Errorf("got %q", buf.String())
		}
		if f.disconnectCalls != 0 {
			t.Errorf("Disconnect must not be called when not connected")
		}
	})

	t.Run("connected", func(t *testing.T) {
		f := &fakeCharger{connected: true}
		ctx, buf := newCtx(f, cfg16())
		handleDisconnect(ctx, nil)
		if !strings.Contains(buf.String(), "Disconnected from server") {
			t.Errorf("got %q", buf.String())
		}
		if f.disconnectCalls != 1 || f.connected {
			t.Errorf("Disconnect not applied: calls=%d connected=%v", f.disconnectCalls, f.connected)
		}
	})
}
