package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestHandleSoc(t *testing.T) {
	t.Run("usage when no arg", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handleSoc(ctx, nil)
		if !strings.Contains(buf.String(), "Usage: soc <0-100>") {
			t.Errorf("got %q", buf.String())
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		f := &fakeCharger{soc: 20}
		ctx, buf := newCtx(f, cfg16())
		handleSoc(ctx, []string{"abc"})
		if !strings.Contains(buf.String(), "Error: invalid SOC value: abc") {
			t.Errorf("got %q", buf.String())
		}
		if f.soc != 20 {
			t.Errorf("SOC must be unchanged on invalid input, got %v", f.soc)
		}
	})

	t.Run("success", func(t *testing.T) {
		f := &fakeCharger{}
		ctx, buf := newCtx(f, cfg16())
		handleSoc(ctx, []string{"55"})
		if !strings.Contains(buf.String(), "SOC set to: 55.0%") {
			t.Errorf("got %q", buf.String())
		}
		if f.soc != 55 {
			t.Errorf("SOC not set, got %v", f.soc)
		}
	})

	t.Run("range error", func(t *testing.T) {
		f := &fakeCharger{setSOCErr: errors.New("SOC must be between 0 and 100")}
		ctx, buf := newCtx(f, cfg16())
		handleSoc(ctx, []string{"150"})
		if !strings.Contains(buf.String(), "Error: SOC must be between 0 and 100") {
			t.Errorf("got %q", buf.String())
		}
	})
}
