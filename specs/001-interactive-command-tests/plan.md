# Implementation Plan: Interactive Command Test Coverage

**Branch**: `001-interactive-command-tests` | **Date**: 2026-07-21 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-interactive-command-tests/spec.md`

## Summary

Refactor the monolithic `interactiveLoop` `switch` in `main.go` into a `cli` package where
every interactive command is a dedicated handler function registered in a command table.
Handlers depend on a `cli.Charger` **interface** (satisfied by `*charger.Charger`) and write
to an injected `io.Writer`, so each handler can be exercised in isolation against a fake
charger with captured output — no real stdin, no network, no live OCPP server. Every handler
gets its own Go test function (`TestHandleConnect`, `TestHandleSoc`, …), plus tests for the
shared parser (whitespace/empty/case-insensitivity/unknown) and a coverage guard that fails
if any registered command lacks a handler entry. Observable command output and state effects
are preserved verbatim by the refactor.

## Technical Context

**Language/Version**: Go 1.25.5 (Golang only, per user directive)

**Primary Dependencies**: Standard library only for the CLI/test layer (`bufio`, `strings`,
`io`, `fmt`, `testing`). Existing modules (`wlgows`, `yaml.v3`) remain used by the `charger`
and `config` packages; the new `cli` package adds no third-party dependencies.

**Storage**: N/A

**Testing**: Go's built-in `testing` package (`go test ./...`). Table-driven tests with an
in-package fake implementing `cli.Charger`; output captured via `bytes.Buffer`.

**Target Platform**: Cross-platform CLI binary (developed on darwin/arm64)

**Project Type**: Single-project CLI application (Go module)

**Performance Goals**: Full interactive-command test suite completes in < 10 s with zero
network access (SC-003); zero flakiness across ≥ 10 runs (SC-004).

**Constraints**: The refactor MUST preserve current observable behavior — the exact strings
each command prints and the charger state each command mutates. No OCPP protocol behavior
changes. No new external dependencies.

**Scale/Scope**: 16 command verbs + shared parsing + unknown/empty paths. ~1 handler file,
~1 registry file, ~1 interface/context file, and a parallel `_test.go` file per unit.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Impact | Verdict |
|-----------|--------|---------|
| I. Protocol Fidelity | Refactor is behavior-preserving; no OCPP message changes. Tests pin protocol-command output. | PASS |
| II. CLI-First Design | CLI commands unchanged in name, syntax, and I/O. | PASS |
| III. Predictable Behavior | Tests are deterministic; fake charger removes network/timing nondeterminism. | PASS |
| IV. Configuration over Code | Version-dependent output still driven by `config`; tests cover 1.6 and 2.0.1. | PASS |
| V. Clear Error Messages | Existing error/usage strings preserved verbatim and asserted by tests. | PASS |
| VI. CLI Output | Output content unchanged; only its destination becomes an injected `io.Writer`. | PASS |
| VII. Testability | Directly fulfills this principle: every command becomes independently testable. | PASS |
| VIII. Extensibility | Registry map makes adding a command (and its required test) a localized change. | PASS |
| IX. Maintainability | Business logic (charger) already separate; this extracts CLI parsing/dispatch into its own `cli` package. | PASS |

**Decision Order alignment**: OCPP correctness is protected (behavior-preserving + tested);
UX unchanged; simplicity favored (stdlib only, one small package); extensibility improved via
the registry. No violations — Complexity Tracking not required.

**Post-Design re-check**: Design introduces only an interface seam, a context struct, a
handler-per-command table, and a fake. No principle is compromised. GATE PASS (re-confirmed).

## Project Structure

### Documentation (this feature)

```text
specs/001-interactive-command-tests/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── commands.md      # Phase 1 output: per-command CLI contract
├── checklists/
│   └── requirements.md  # From /speckit-specify
└── tasks.md             # Phase 2 output (/speckit-tasks — NOT created here)
```

### Source Code (repository root)

```text
main.go                  # SLIMMED: load config, create charger, wire signals,
                         #   call cli.Run(sim, cfg, os.Stdin, os.Stdout)

cli/
├── charger.go           # Charger interface (methods the handlers depend on)
├── context.go           # CommandContext { Charger, Config, Out io.Writer }
├── registry.go          # command table: map[string]Handler (+ aliases), Dispatch, Run loop
├── handlers.go          # one handler func per command (handleConnect, handleSoc, …)
├── help.go              # handleHelp + printHelp (moved from main.go)
├── fake_charger_test.go # fakeCharger implementing cli.Charger for tests
├── registry_test.go     # parser tests + coverage guard (every command has a handler)
├── handlers_test.go     # TestHandle<Command> — one test func per handler
└── help_test.go         # help/info version-dependent output tests

charger/                 # UNCHANGED (business logic; already satisfies cli.Charger)
config/                  # UNCHANGED
ocpp/                    # UNCHANGED
```

**Structure Decision**: Single Go module. Introduce one new package, `cli/`, that owns command
parsing, the command registry, and per-command handlers. `main.go` is reduced to wiring
(config load, charger construction, signal handling) and delegates the interactive loop to
`cli.Run`. This satisfies Constitution Principle IX (separate CLI parsing from business logic)
and makes each handler unit-testable without a `main`-package test. `*charger.Charger` already
implements every method the handlers need, so it satisfies `cli.Charger` with no changes to the
`charger` package.

## Complexity Tracking

> No constitution violations. Section intentionally left empty.
