package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandleConnect(t *testing.T) {
	t.Run("already connected", func(t *testing.T) {
		f := &fakeCharger{connected: true}
		ctx, buf := newCtx(f, cfg16())
		handleConnect(ctx, nil)
		if !strings.Contains(buf.String(), "Already connected") {
			t.Errorf("got %q", buf.String())
		}
		if f.connectCalls != 0 {
			t.Errorf("Connect must not be called when already connected, calls=%d", f.connectCalls)
		}
	})

	t.Run("success chain", func(t *testing.T) {
		f := &fakeCharger{status: "Available"}
		ctx, buf := newCtx(f, cfg16())
		handleConnect(ctx, nil)
		out := buf.String()
		if !strings.Contains(out, "Connected to server") {
			t.Errorf("missing connected message: %q", out)
		}
		if f.bootCalls != 1 {
			t.Errorf("BootNotification not sent, bootCalls=%d", f.bootCalls)
		}
		if !f.statusNotifSet || f.statusNotifArg != "Available" {
			t.Errorf("StatusNotification not sent with current status: set=%v arg=%q", f.statusNotifSet, f.statusNotifArg)
		}
	})

	t.Run("connect error", func(t *testing.T) {
		f := &fakeCharger{connectErr: errors.New("dial fail")}
		ctx, buf := newCtx(f, cfg16())
		handleConnect(ctx, nil)
		if !strings.Contains(buf.String(), "Error: dial fail") {
			t.Errorf("got %q", buf.String())
		}
		if f.bootCalls != 0 {
			t.Errorf("BootNotification must not be sent after a connect failure")
		}
	})
}
