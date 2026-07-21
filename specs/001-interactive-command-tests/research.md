# Phase 0 Research: Interactive Command Test Coverage

All Technical Context items are resolved; there are no remaining NEEDS CLARIFICATION. The
user directives ("Golang only", "each command has its own handler", "each handler has its own
test func") fully constrain the approach. The decisions below record how the spec's
requirements map onto a Go design.

## Decision 1: Test seam via a `cli.Charger` interface

- **Decision**: Define an interface `cli.Charger` listing exactly the methods the command
  handlers call. Handlers depend on the interface, not on `*charger.Charger`. Production wires
  the real `*charger.Charger`; tests wire a `fakeCharger`.
- **Rationale**: `charger.Charger.Connect()` dials a real WebSocket. Depending on an interface
  is the idiomatic Go seam that lets connection- and protocol-dependent commands run
  deterministically offline (FR-001, SC-003). `*charger.Charger` already implements every
  method, so no change to the `charger` package is required.
- **Interface members** (confirmed against the codebase):
  `IsConnected() bool`, `Connect() error`, `Disconnect()`, `BootNotification() error`,
  `StatusNotification(status string) error`, `GetStatus() string`, `SetStatus(status string) error`,
  `Plugin() error`, `Unplug() error`, `StartTransaction(idTag string) error`,
  `StopTransaction(reason string) error`, `MeterValues() error`,
  `SetLicensePlateAndSend(plate string) error`, `GetLicensePlate() string`,
  `SetSOC(soc float64) error`, `GetSOC() float64`, `SetCurrent(a float64) error`,
  `GetCurrent() float64`, `SetPower(w float64) error`, `GetPower() float64`, `IsCharging() bool`.
- **Alternatives considered**: (a) `httptest`-based fake OCPP server — rejected: slower, more
  fragile, exercises the transport rather than the command layer under test. (b) Build tags to
  swap implementations — rejected: heavier than an interface and harder to read.

## Decision 2: One handler per command with a uniform signature

- **Decision**: Each command is a function `func(ctx *CommandContext, args []string)`. A
  `registry map[string]Handler` maps the lower-cased command word to its handler; aliases
  (`exit` → same handler as `quit`) are separate map entries.
- **Rationale**: Directly satisfies the user's "each command has its own handler". A uniform
  signature makes handlers trivially testable and the registry makes dispatch and the coverage
  guard (FR-011) simple. `args` excludes the command word (matches today's `parts[1:]` usage).
- **Alternatives considered**: Keep the `switch` and test through it — rejected: the user
  explicitly asked for per-command handlers, and a `switch` is harder to enumerate for the
  coverage guard.

## Decision 3: Output via injected `io.Writer`

- **Decision**: `CommandContext` carries `Out io.Writer`. Handlers use `fmt.Fprintf(ctx.Out, …)`
  / `fmt.Fprintln` instead of `fmt.Printf` / `fmt.Println`. Production passes `os.Stdout`; tests
  pass a `*bytes.Buffer`.
- **Rationale**: Lets tests capture and assert exact command output (FR-007) while preserving
  the current wording verbatim (Principle V/VI). Only the *destination* changes.
- **Note**: Diagnostic lines emitted by `charger` internals via the standard `log` package
  (e.g. "Connecting to …") are out of scope for handler assertions and remain on the logger.

## Decision 4: Table-driven Go tests, one `Test` func per handler

- **Decision**: A dedicated test function per handler (`TestHandleConnect`, `TestHandleSoc`, …),
  each using table-driven sub-cases (`t.Run`) for its success path, guard/error path, and
  usage/invalid-input path. A `fakeCharger` records calls and returns programmable errors.
- **Rationale**: Satisfies "each handler has its own test func" while keeping success/error
  coverage (FR-006) compact. Sub-tests give per-scenario failure names (SC-005).
- **Alternatives considered**: One giant table across all commands — rejected: violates the
  user's per-handler-test-func requirement and produces poor failure localization.

## Decision 5: Parser extraction and shared-behavior tests

- **Decision**: Extract line handling into a pure `parseCommand(line string) (cmd string, args []string, ok bool)`
  that trims, splits on whitespace, lower-cases the verb, and reports `ok=false` for empty
  input. `Dispatch` looks up the verb and falls back to the unknown-command message. The blocking
  read loop (`Run`) reads from an `io.Reader` and returns on EOF (FR-002).
- **Rationale**: Isolates the rules shared by every command (FR-005) so they are tested once,
  and makes `Run` terminate on input exhaustion so tests can feed a finite script.
- **Alternatives considered**: Parsing inside each handler — rejected: duplicates logic
  (Principle IX) and can't be tested independently.

## Decision 6: Coverage guard for "every command is tested"

- **Decision**: Maintain a canonical list of expected command verbs. A test asserts the
  registry contains a handler for each, and a documented convention (one `TestHandle<Command>`
  per handler) makes an untested command visible. Adding a command without a registry entry
  fails the guard.
- **Rationale**: Fulfills FR-011 and SC-006 — an untested/unregistered command is detectable
  from the suite itself.

## Decision 7: Version-dependent output tested for both versions

- **Decision**: `help` and `status`-with-no-arg tests run twice, once with a `config.Config`
  where `OCPPVersion == "1.6"` and once `"2.0.1"`, asserting the correct valid-status set.
- **Rationale**: Satisfies FR-008 / US4-AC4 / US5-AC1. `config.Config` is a plain struct that
  can be constructed directly in tests (no file load needed).

## Decision 8: Package placement

- **Decision**: New package `cli` at the module root; `main` shrinks to wiring and calls
  `cli.Run`. Tests live in the `cli` package (`package cli`) alongside the code.
- **Rationale**: Testing a `main` package is awkward; a dedicated package makes handlers and the
  registry directly importable by tests and aligns with Principle IX. "Golang only" honored —
  standard library, standard `go test`, no new tooling.
