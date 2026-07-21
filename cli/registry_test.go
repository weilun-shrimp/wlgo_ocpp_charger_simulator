package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// T009 [US1]: Run must return when its input reader is exhausted (EOF).
func TestRun_TerminatesOnEOF(t *testing.T) {
	f := &fakeCharger{status: "Available"}
	in := strings.NewReader("info\nquit\n")
	out := &bytes.Buffer{}

	done := make(chan struct{})
	go func() {
		Run(f, cfg16(), in, out)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return on EOF")
	}

	if !strings.Contains(out.String(), "Use Ctrl+C to exit") {
		t.Errorf("expected commands to have run before EOF, got: %q", out.String())
	}
}

// T010 [US1]: Dispatch routes to the handler, captures output into the injected
// writer, and touches only the fake charger — no real network.
func TestDispatch_CapturesOutputNoNetwork(t *testing.T) {
	f := &fakeCharger{status: "Available"}
	ctx, buf := newCtx(f, cfg16())

	Dispatch(ctx, "info", nil)
	if !strings.Contains(buf.String(), "Status: Available") {
		t.Errorf("expected captured info output, got %q", buf.String())
	}

	// On a real charger "connect" dials a WebSocket; here it only exercises the fake.
	Dispatch(ctx, "connect", nil)
	if f.connectCalls != 1 || !f.connected {
		t.Errorf("expected fake connect to be used, calls=%d connected=%v", f.connectCalls, f.connected)
	}
}

// T011 [US2]: parseCommand trims, skips empty/whitespace, lower-cases the verb,
// and splits arguments preserving their case.
func TestParseCommand(t *testing.T) {
	cases := []struct {
		name     string
		line     string
		wantCmd  string
		wantArgs []string
		wantOK   bool
	}{
		{"empty", "", "", nil, false},
		{"whitespace only", "   \t  ", "", nil, false},
		{"lowercases verb", "CONNECT", "connect", nil, true},
		{"mixed case verb", "Info", "info", nil, true},
		{"single arg", "soc 55", "soc", []string{"55"}, true},
		{"arg case preserved", "plate ABC123", "plate", []string{"ABC123"}, true},
		{"trims and splits", "  status   Charging  ", "status", []string{"Charging"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, args, ok := parseCommand(tc.line)
			if cmd != tc.wantCmd || ok != tc.wantOK {
				t.Fatalf("parseCommand(%q) = (%q, %v, %v); want (%q, _, %v)",
					tc.line, cmd, args, ok, tc.wantCmd, tc.wantOK)
			}
			if !equalArgs(args, tc.wantArgs) {
				t.Errorf("args = %v; want %v", args, tc.wantArgs)
			}
		})
	}
}

// T012 [US2]: an unrecognized verb prints the unknown-command message, and
// empty input is ignored (no command dispatched).
func TestDispatch_UnknownAndEmpty(t *testing.T) {
	f := &fakeCharger{}
	ctx, buf := newCtx(f, cfg16())

	Dispatch(ctx, "frobnicate", nil)
	if !strings.Contains(buf.String(), "Unknown command: frobnicate. Type 'help' for available commands.") {
		t.Errorf("unknown-command message missing, got %q", buf.String())
	}

	if _, _, ok := parseCommand("   "); ok {
		t.Error("expected empty/whitespace input to be ignored (ok=false)")
	}
}
