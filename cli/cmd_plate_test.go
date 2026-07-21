package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandlePlate(t *testing.T) {
	t.Run("usage when no arg", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handlePlate(ctx, nil)
		if !strings.Contains(buf.String(), "Usage: plate <license_plate>") {
			t.Errorf("got %q", buf.String())
		}
		if f.lastPlate != "" {
			t.Errorf("SetLicensePlateAndSend must not be called without an argument")
		}
	})

	t.Run("success", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handlePlate(ctx, []string{"ABC123"})
		if !strings.Contains(buf.String(), "License plate set: ABC123") {
			t.Errorf("got %q", buf.String())
		}
		if f.licensePlate != "ABC123" {
			t.Errorf("plate not set, got %q", f.licensePlate)
		}
	})

	t.Run("error", func(t *testing.T) {
		f := &fakeCharger{plateErr: errors.New("send failed")}
		ctx, buf := newCtx(f, cfg16())
		handlePlate(ctx, []string{"ABC123"})
		if !strings.Contains(buf.String(), "Error: send failed") {
			t.Errorf("got %q", buf.String())
		}
	})
}
