package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandleStop(t *testing.T) {
	t.Run("default reason", func(t *testing.T) {
		f := &fakeCharger{charging: true}
		ctx, buf := newCtx(f, cfg16())
		handleStop(ctx, nil)
		if !strings.Contains(buf.String(), "Transaction stopped") {
			t.Errorf("got %q", buf.String())
		}
		if f.lastStopReason != "Local" {
			t.Errorf("expected default reason Local, got %q", f.lastStopReason)
		}
	})

	t.Run("explicit reason", func(t *testing.T) {
		f := &fakeCharger{charging: true}
		ctx, _ := newCtx(f, cfg16())
		handleStop(ctx, []string{"Remote"})
		if f.lastStopReason != "Remote" {
			t.Errorf("expected reason Remote, got %q", f.lastStopReason)
		}
	})

	t.Run("error", func(t *testing.T) {
		f := &fakeCharger{stopErr: errors.New("no transaction")}
		ctx, buf := newCtx(f, cfg16())
		handleStop(ctx, nil)
		if !strings.Contains(buf.String(), "Error: no transaction") {
			t.Errorf("got %q", buf.String())
		}
	})
}
