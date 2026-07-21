package cli

import (
	"strings"
	"testing"
)

func TestHandleInfo(t *testing.T) {
	t.Run("without plate", func(t *testing.T) {
		f := &fakeCharger{connected: true, status: "Charging", charging: true, current: 10, power: 2300, soc: 42}
		ctx, buf := newCtx(f, cfg16())
		handleInfo(ctx, nil)
		out := buf.String()
		for _, want := range []string{
			"Connected: true",
			"Status: Charging",
			"Charging: true",
			"Voltage: 230.0 V",
			"Current: 10.0 A",
			"Power: 2300.0 W",
			"SOC: 42.0%",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q in output: %q", want, out)
			}
		}
		if strings.Contains(out, "License Plate") {
			t.Errorf("license-plate line must be absent when no plate is set: %q", out)
		}
	})

	t.Run("with plate", func(t *testing.T) {
		f := &fakeCharger{licensePlate: "ABC123"}
		ctx, buf := newCtx(f, cfg16())
		handleInfo(ctx, nil)
		if !strings.Contains(buf.String(), "License Plate: ABC123") {
			t.Errorf("expected license-plate line, got %q", buf.String())
		}
	})
}
