# Quickstart & Validation: Interactive Command Test Coverage

This guide validates the feature end-to-end. It assumes the `cli` package and its tests exist
(produced during implementation). No network or external OCPP server is required.

## Prerequisites

- Go 1.25+ installed (`go version`)
- Repository checked out on branch `001-interactive-command-tests`
- From the repository root

## Run the interactive-command test suite

```bash
# Run just the new CLI package tests, verbose (shows one line per handler test)
go test ./cli/ -v

# Run everything and confirm nothing else broke
go test ./...
```

**Expected**: All tests pass. Verbose output shows a `Test<...>` entry for each handler —
`TestHandleHelp`, `TestHandleConnect`, `TestHandleDisconnect`, `TestHandlePlugin`,
`TestHandleUnplug`, `TestHandleStatus`, `TestHandleStart`, `TestHandleStop`, `TestHandleMeter`,
`TestHandlePlate`, `TestHandleSoc`, `TestHandleCurrent`, `TestHandlePower`, `TestHandleInfo`,
`TestHandleQuit` — plus `TestParseCommand` and `TestRegistryCoversAllCommands`.

## Validate the acceptance criteria

| Spec item | How to validate |
|-----------|-----------------|
| SC-001 (every command tested) | `go test ./cli/ -v` shows a `TestHandle*` per command; `TestRegistryCoversAllCommands` passes |
| SC-002 (both paths) | Handler tests include success and guard/error sub-cases (`go test ./cli/ -run TestHandleConnect -v`) |
| SC-003 (< 10 s, no network) | Disconnect from any network, then `time go test ./cli/` completes quickly |
| SC-004 (no flakiness) | `go test ./cli/ -count=10` passes every iteration |
| SC-006 (coverage detectable) | Temporarily add an unregistered command verb → `TestRegistryCoversAllCommands` fails |

## Prove the suite catches regressions (SC-005)

```bash
# Example: change a handler's success message, e.g. "Connected to server" -> "Connected!"
# Then:
go test ./cli/ -run TestHandleConnect
# Expected: the test FAILS, demonstrating the suite pins observable output.
# Revert the change afterward.
```

## Confirm real behavior is unchanged (behavior-preserving refactor)

```bash
# Build and run the simulator as before; commands behave identically.
go build ./...
./wlgo_ocpp_charger_simulator -config config.yaml
# At the prompt, try: help, info, soc 55, status, quit
```

**Expected**: Identical output and behavior to before the refactor — only the internal wiring
(handlers + injected writer) changed.

## References

- Per-command behavior: [contracts/commands.md](./contracts/commands.md)
- Types and seams: [data-model.md](./data-model.md)
- Design decisions: [research.md](./research.md)
