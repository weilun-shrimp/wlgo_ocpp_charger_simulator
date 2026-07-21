package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandleMeter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handleMeter(ctx, nil)
		if !strings.Contains(buf.String(), "MeterValues updated") {
			t.Errorf("got %q", buf.String())
		}
		if f.meterCalls != 1 {
			t.Errorf("MeterValues not called, meterCalls=%d", f.meterCalls)
		}
	})

	t.Run("error", func(t *testing.T) {
		f := &fakeCharger{meterErr: errors.New("send failed")}
		ctx, buf := newCtx(f, cfg16())
		handleMeter(ctx, nil)
		if !strings.Contains(buf.String(), "Error: send failed") {
			t.Errorf("got %q", buf.String())
		}
	})
}
