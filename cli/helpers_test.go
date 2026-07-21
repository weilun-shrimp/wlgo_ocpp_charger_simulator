package cli

import (
	"bytes"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/config"
)

// newCtx builds a CommandContext wired to the given fake charger and config,
// with output captured into the returned buffer.
func newCtx(f *fakeCharger, cfg *config.Config) (*CommandContext, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return &CommandContext{Charger: f, Config: cfg, Out: buf}, buf
}

// cfg16 returns a minimal OCPP 1.6 config for tests.
func cfg16() *config.Config {
	return &config.Config{OCPPVersion: "1.6", Voltage: 230, MaxCurrent: 32, MaxPower: 7360}
}

// cfg201 returns a minimal OCPP 2.0.1 config for tests.
func cfg201() *config.Config {
	return &config.Config{OCPPVersion: "2.0.1", Voltage: 230, MaxCurrent: 32, MaxPower: 7360}
}

// equalArgs compares two argument slices, treating nil and empty as equal.
func equalArgs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
