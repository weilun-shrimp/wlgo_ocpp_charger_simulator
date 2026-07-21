# Phase 1 Data Model: Interactive Command Test Coverage

This feature is structural (a testable CLI layer), so the "entities" are Go types introduced
in the new `cli` package plus the test fixtures. No persistent data or storage is involved.

## Type: `Charger` (interface)

The seam handlers depend on. Implemented in production by `*charger.Charger` (unchanged), and
in tests by `fakeCharger`.

| Member | Signature | Used by commands |
|--------|-----------|------------------|
| IsConnected | `() bool` | connect, disconnect, info |
| Connect | `() error` | connect |
| Disconnect | `()` | disconnect |
| BootNotification | `() error` | connect |
| StatusNotification | `(status string) error` | connect |
| GetStatus | `() string` | connect, status, info |
| SetStatus | `(status string) error` | status |
| Plugin | `() error` | plugin |
| Unplug | `() error` | unplug |
| StartTransaction | `(idTag string) error` | start |
| StopTransaction | `(reason string) error` | stop |
| MeterValues | `() error` | meter |
| SetLicensePlateAndSend | `(plate string) error` | plate |
| GetLicensePlate | `() string` | info |
| SetSOC | `(soc float64) error` | soc |
| GetSOC | `() float64` | info |
| SetCurrent | `(a float64) error` | current |
| GetCurrent | `() float64` | current, info |
| SetPower | `(w float64) error` | power |
| GetPower | `() float64` | power, info |
| IsCharging | `() bool` | info |

## Type: `CommandContext`

Passed to every handler; the only state a handler may touch.

| Field | Type | Purpose |
|-------|------|---------|
| Charger | `Charger` | The (real or fake) charger the command acts on |
| Config | `*config.Config` | Version + limits used for usage text and version-dependent output |
| Out | `io.Writer` | Destination for all command output (stdout in prod, buffer in tests) |

## Type: `Handler`

```
type Handler func(ctx *CommandContext, args []string)
```

- `args` are the whitespace-split tokens **after** the command word.
- A handler writes all user-facing text to `ctx.Out` and never to `os.Stdout` directly.

## Type: `registry`

| Aspect | Definition |
|--------|------------|
| Shape | `map[string]Handler` keyed by lower-cased command verb |
| Aliases | `quit` and `exit` map to the same handler entry |
| Lookup | `Dispatch(ctx, cmd, args)`: found → call handler; not found → print unknown-command message |
| Canonical verbs | help, connect, disconnect, plugin, unplug, status, start, stop, meter, plate, soc, current, power, info, quit, exit |

## Function: `parseCommand`

```
parseCommand(line string) (cmd string, args []string, ok bool)
```

| Rule | Behavior |
|------|----------|
| Trim | Leading/trailing whitespace removed |
| Empty | Empty or whitespace-only → `ok == false` (no command runs) |
| Split | Remaining tokens split on whitespace |
| Verb case | `cmd` is lower-cased; `args` preserve original case |

## Function: `Run` (loop)

```
Run(c Charger, cfg *config.Config, in io.Reader, out io.Writer)
```

- Reads lines from `in`, prints the `> ` prompt to `out`, dispatches each parsed command.
- **Returns on EOF / input exhaustion** (FR-002) so tests can feed a finite script.

## Test fixture: `fakeCharger`

In-package test double implementing `Charger`.

| Concern | Behavior |
|---------|----------|
| State | Holds `connected`, `status`, `soc`, `current`, `power`, `charging`, `licensePlate` fields |
| Programmable errors | Per-method error hooks (e.g. `connectErr`, `pluginErr`, `startErr`) default nil |
| Call recording | Records invocations/arguments (e.g. `lastStopReason`, `bootCalled`) for assertions |
| Determinism | Pure in-memory; no goroutines, no network, no timers |

## Test fixture: command test case

Conceptual shape used by table-driven handler tests.

| Field | Meaning |
|-------|---------|
| name | Sub-test name (scenario) |
| args | Tokens passed to the handler |
| setup | Initial `fakeCharger`/`Config` state (e.g. pre-connected, `Preparing`, `OCPPVersion`) |
| wantOut | Expected substring(s) of captured output |
| wantState | Expected `fakeCharger` state / recorded call after the handler runs |

## State transitions (unchanged; asserted, not modified)

The handlers do not add transitions — they invoke existing `charger` behavior. Tests assert the
observable results the current code already produces, e.g.:

- `connect` from disconnected → `Connect` + `BootNotification` + `StatusNotification` attempted.
- `connect` while connected → "Already connected"; no second `Connect`.
- `disconnect` while disconnected → "Not connected".
- `status`/`soc`/`current`/`power`/`plate` with no arg → that command's usage line.
- `soc`/`current`/`power` with non-numeric arg → invalid-value message; no state change.
- `stop` with no arg → default reason `Local`.
