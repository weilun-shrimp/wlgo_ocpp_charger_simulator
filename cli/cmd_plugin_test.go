package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandlePlugin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handlePlugin(ctx, nil)
		if !strings.Contains(buf.String(), "Car plugged in (Preparing)") {
			t.Errorf("got %q", buf.String())
		}
	})

	t.Run("error", func(t *testing.T) {
		f := &fakeCharger{pluginErr: errors.New("must be Available")}
		ctx, buf := newCtx(f, cfg16())
		handlePlugin(ctx, nil)
		if !strings.Contains(buf.String(), "Error: must be Available") {
			t.Errorf("got %q", buf.String())
		}
	})
}
