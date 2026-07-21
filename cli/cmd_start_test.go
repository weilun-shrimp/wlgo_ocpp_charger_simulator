package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandleStart(t *testing.T) {
	t.Run("usage when no idTag", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handleStart(ctx, nil)
		if !strings.Contains(buf.String(), "Usage: start <idTag>") {
			t.Errorf("got %q", buf.String())
		}
		if f.lastStartIDTag != "" {
			t.Errorf("StartTransaction must not be called without an idTag")
		}
	})

	t.Run("success", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handleStart(ctx, []string{"TAG1"})
		if !strings.Contains(buf.String(), "Transaction started") {
			t.Errorf("got %q", buf.String())
		}
		if f.lastStartIDTag != "TAG1" || !f.charging {
			t.Errorf("start not applied: idTag=%q charging=%v", f.lastStartIDTag, f.charging)
		}
	})

	t.Run("error", func(t *testing.T) {
		f := &fakeCharger{startErr: errors.New("not preparing")}
		ctx, buf := newCtx(f, cfg16())
		handleStart(ctx, []string{"TAG1"})
		if !strings.Contains(buf.String(), "Error: not preparing") {
			t.Errorf("got %q", buf.String())
		}
	})
}
