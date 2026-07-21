---
description: "Task list for Interactive Command Test Coverage"
---

# Tasks: Interactive Command Test Coverage

**Input**: Design documents from `/specs/001-interactive-command-tests/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/commands.md

**Tests**: This feature IS a testing feature — test tasks are REQUIRED (one dedicated test
function per handler, per the spec and user directive "each handler can have its own task").

**Organization**: Tasks are grouped by user story. Per the user directive, every command handler
gets its own implementation task and its own test task, each in its own file so they run in
parallel.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: Which user story this task belongs to (US1–US5)
- Every task lists an exact file path

## Path Conventions

- Single Go module at repository root. New package `cli/`. Handlers and their tests live in
  per-command files inside `cli/`. Existing `charger/`, `config/`, `ocpp/` packages are unchanged.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the new package that will hold the CLI command layer.

- [X] T001 Create the `cli` package with a package doc file in `cli/doc.go` (`// Package cli implements the interactive command layer.` + `package cli`)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The testable dispatch harness every command depends on. Realizes User Story 1's
infrastructure. **No handler or handler test can begin until this phase is complete.**

**⚠️ CRITICAL**: Blocks all user story phases.

- [X] T002 [P] Define the `Charger` interface (all methods from data-model.md) in `cli/charger.go`; confirm `*charger.Charger` satisfies it (no changes to the `charger` package)
- [X] T003 [P] Define `CommandContext{ Charger, Config *config.Config, Out io.Writer }` in `cli/context.go`
- [X] T004 Define `Handler func(*CommandContext, []string)`, the `registry map[string]Handler`, a `register(name string, h Handler)` helper, and `Dispatch(ctx, cmd, args)` with the unknown-command fallback (`Unknown command: <cmd>. Type 'help' for available commands.`) in `cli/registry.go`
- [X] T005 Implement `parseCommand(line) (cmd string, args []string, ok bool)` (trim, split on whitespace, lower-case verb, `ok=false` on empty/whitespace) in `cli/registry.go` (same file as T004)
- [X] T006 Implement `Run(c Charger, cfg *config.Config, in io.Reader, out io.Writer)` — print `> ` prompt, read lines, dispatch, and RETURN on EOF/input exhaustion — in `cli/registry.go` (same file as T004, T005)
- [X] T007 Slim `main.go` to load config, create the charger, wire signal handling, and call `cli.Run(sim, cfg, os.Stdin, os.Stdout)`; remove the old `interactiveLoop`/`printHelp` from `main.go`
- [X] T008 [P] Create the `fakeCharger` test double implementing `cli.Charger` (in-memory state + per-method programmable error hooks + call recording, no network/goroutines/timers) in `cli/fake_charger_test.go`

**Checkpoint**: Harness compiles and the fake exists — command handlers and their tests can now be built in parallel.

---

## Phase 3: User Story 1 - Testable command dispatch harness (Priority: P1)

**Goal**: Prove a command can be driven with controlled input, its output captured, run against a
fake (no network), and that the loop returns on input exhaustion.

**Independent Test**: Feed a finite input script into `Run`/`Dispatch` with a `bytes.Buffer` and a
`fakeCharger`; assert captured output and that `Run` returns.

- [X] T009 [US1] `TestRun_TerminatesOnEOF` — `Run` returns when the input reader is exhausted — in `cli/registry_test.go`
- [X] T010 [US1] `TestDispatch_CapturesOutputNoNetwork` — dispatching a known command writes to the injected buffer and touches only the `fakeCharger` (no real connection) — in `cli/registry_test.go` (same file as T009)

**Checkpoint**: Harness verified — every subsequent handler test builds on this pattern.

---

## Phase 4: User Story 2 - Input parsing behavior is verified (Priority: P1)

**Goal**: Lock down the shared parsing rules that front every command.

**Independent Test**: Submit blank/whitespace/mixed-case/unknown/extra-arg inputs and assert each is
handled per contract without invoking a command action.

- [X] T011 [US2] `TestParseCommand` — table cases for whitespace trim, empty/whitespace → `ok=false`, case-insensitive verb, argument splitting — in `cli/registry_test.go`
- [X] T012 [US2] `TestDispatch_UnknownAndEmpty` — unknown verb prints the unknown-command message; empty input runs no command — in `cli/registry_test.go` (same file as T011)

**Checkpoint**: Parsing rules pinned independently of any single command.

---

## Phase 5: User Story 3 - Session lifecycle commands (Priority: P1)

**Goal**: Each lifecycle command has its own handler and its own test (success + guard/error paths).

**Independent Test**: Run each command from a known `fakeCharger` state; assert printed message and
resulting state; then run its guard/error path.

### Handlers (one per file, all parallel)

- [X] T013 [P] [US3] Implement `handleConnect` (already-connected guard → `Already connected`; success → `Connected to server` + `BootNotification` + `StatusNotification`; error branches) and register `connect` in `cli/cmd_connect.go`
- [X] T014 [P] [US3] Implement `handleDisconnect` (not-connected guard → `Not connected`; success → `Disconnected from server`) and register `disconnect` in `cli/cmd_disconnect.go`
- [X] T015 [P] [US3] Implement `handlePlugin` (success → `Car plugged in (Preparing)`; error → `Error: <err>`) and register `plugin` in `cli/cmd_plugin.go`
- [X] T016 [P] [US3] Implement `handleUnplug` (success → `Car unplugged (Available)`; error → `Error: <err>`) and register `unplug` in `cli/cmd_unplug.go`
- [X] T017 [P] [US3] Implement `handleStart` (no arg → `Usage: start <idTag>`; success → `Transaction started`; error → `Error: <err>`) and register `start` in `cli/cmd_start.go`
- [X] T018 [P] [US3] Implement `handleStop` (default reason `Local`; success → `Transaction stopped`; error → `Error: <err>`) and register `stop` in `cli/cmd_stop.go`
- [X] T019 [P] [US3] Implement `handleMeter` (success → `MeterValues updated`; error → `Error: <err>`) and register `meter` in `cli/cmd_meter.go`

### Tests (one per handler, all parallel; each depends on its handler + T008)

- [X] T020 [P] [US3] `TestHandleConnect` (guard, success chain, connect error) in `cli/cmd_connect_test.go`
- [X] T021 [P] [US3] `TestHandleDisconnect` (both paths) in `cli/cmd_disconnect_test.go`
- [X] T022 [P] [US3] `TestHandlePlugin` (success + error) in `cli/cmd_plugin_test.go`
- [X] T023 [P] [US3] `TestHandleUnplug` (success + error) in `cli/cmd_unplug_test.go`
- [X] T024 [P] [US3] `TestHandleStart` (usage, success, error) in `cli/cmd_start_test.go`
- [X] T025 [P] [US3] `TestHandleStop` (default reason, explicit reason, error) in `cli/cmd_stop_test.go`
- [X] T026 [P] [US3] `TestHandleMeter` (success + error) in `cli/cmd_meter_test.go`

**Checkpoint**: All core OCPP lifecycle commands independently handled and tested.

---

## Phase 6: User Story 4 - State and parameter commands (Priority: P2)

**Goal**: Each parameter command has its own handler and test (usage / invalid / success paths).

**Independent Test**: For each command submit no-arg, malformed-arg, and valid-arg; assert usage,
invalid-value message, and confirmation + state change respectively.

### Handlers (one per file, all parallel)

- [X] T027 [P] [US4] Implement `handleStatus` (no arg → `Usage: status <status>` + version-appropriate valid-status line; success → `Status updated to: <arg>`; error → `Error: <err>`) and register `status` in `cli/cmd_status.go`
- [X] T028 [P] [US4] Implement `handleSoc` (no arg → `Usage: soc <0-100>`; non-numeric → `Error: invalid SOC value: <arg>`; success → `SOC set to: <v>%`; range error → `Error: <err>`) and register `soc` in `cli/cmd_soc.go`
- [X] T029 [P] [US4] Implement `handleCurrent` (no arg → usage `Usage: current <amperes> (0-<MaxCurrent> A, 0 = SuspendedEVSE)` + `Current: <GetCurrent()> A`; non-numeric → `Error: invalid current value: <arg>`; success → `Current set to: <v> A`; error → `Error: <err>`) and register `current` in `cli/cmd_current.go`
- [X] T030 [P] [US4] Implement `handlePower` (no arg → usage `Usage: power <watts> (0-<MaxPower> W, 0 = SuspendedEVSE)` + `Power: <GetPower()> W`; non-numeric → `Error: invalid power value: <arg>`; success → `Power set to: <v> W`; error → `Error: <err>`) and register `power` in `cli/cmd_power.go`
- [X] T031 [P] [US4] Implement `handlePlate` (no arg → `Usage: plate <license_plate>`; success → `License plate set: <arg>`; error → `Error: <err>`) and register `plate` in `cli/cmd_plate.go`

### Tests (one per handler, all parallel)

- [X] T032 [P] [US4] `TestHandleStatus` (usage for BOTH OCPP 1.6 and 2.0.1, success, error) in `cli/cmd_status_test.go`
- [X] T033 [P] [US4] `TestHandleSoc` (usage, invalid, success, range error) in `cli/cmd_soc_test.go`
- [X] T034 [P] [US4] `TestHandleCurrent` (usage, invalid, success, error) in `cli/cmd_current_test.go`
- [X] T035 [P] [US4] `TestHandlePower` (usage, invalid, success, error) in `cli/cmd_power_test.go`
- [X] T036 [P] [US4] `TestHandlePlate` (usage, success, error) in `cli/cmd_plate_test.go`

**Checkpoint**: All parameter commands independently handled and tested, both OCPP versions covered.

---

## Phase 7: User Story 5 - Informational and meta commands (Priority: P2)

**Goal**: Each read-only/meta command has its own handler and test.

**Independent Test**: Run each and assert output — `help` lists commands + version status set;
`info` reports fields with the plate line only when set; `quit`/`exit` prints the Ctrl+C hint.

### Handlers (one per file, all parallel)

- [X] T037 [P] [US5] Implement `handleHelp` + move `printHelp` (full command list + version-appropriate valid-status line) and register `help` in `cli/cmd_help.go`
- [X] T038 [P] [US5] Implement `handleInfo` (Connected/Status/Charging/Voltage/Current/Power/SOC lines; `License Plate:` line ONLY when `GetLicensePlate() != ""`) and register `info` in `cli/cmd_info.go`
- [X] T039 [P] [US5] Implement `handleQuit` (→ `Use Ctrl+C to exit`) and register BOTH `quit` and `exit` in `cli/cmd_quit.go`

### Tests (one per handler, all parallel)

- [X] T040 [P] [US5] `TestHandleHelp` (BOTH OCPP 1.6 and 2.0.1 status sets) in `cli/cmd_help_test.go`
- [X] T041 [P] [US5] `TestHandleInfo` (plate present vs absent) in `cli/cmd_info_test.go`
- [X] T042 [P] [US5] `TestHandleQuit` (both `quit` and `exit` aliases) in `cli/cmd_quit_test.go`

**Checkpoint**: Every one of the 16 command verbs has its own handler and its own test function.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Coverage guard and full-suite verification.

- [X] T043 [P] `TestRegistryCoversAllCommands` — assert the registry has a handler for every canonical verb (help, connect, disconnect, plugin, unplug, status, start, stop, meter, plate, soc, current, power, info, quit, exit); fails if a command is added without registration — in `cli/coverage_test.go`
- [X] T044 Run `go test ./...` and `go vet ./...`; confirm all pass (SC-001, SC-002) — repo root
- [X] T045 [P] Run `go test ./cli/ -count=10` to confirm zero flakiness (SC-004) and `time go test ./cli/` under 10s with no network (SC-003) — repo root
- [X] T046 [P] `go build ./...` and run the manual smoke steps in `quickstart.md` to confirm behavior-preserving refactor (help, info, soc 55, status, quit) — repo root

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: none — start immediately.
- **Foundational (Phase 2)**: depends on Phase 1 — BLOCKS all user stories.
- **User Stories (Phases 3–7)**: all depend on Phase 2. US1 and US2 tests share `cli/registry_test.go` (sequential within that file). US3/US4/US5 handlers and tests are in per-command files and run fully in parallel.
- **Polish (Phase 8)**: T043 depends on all handlers being registered (Phases 5–7); T044–T046 depend on everything.

### Key task dependencies

- T004 → T005 → T006 (same file `cli/registry.go`, sequential).
- T007 depends on T002 + T006 (`cli.Run` and `Charger` must exist).
- T008 depends on T002 (must implement the `Charger` interface).
- Each handler test (T020–T026, T032–T036, T040–T042) depends on its own handler impl and on T008.
- T043 depends on every handler registration (T013–T019, T027–T031, T037–T039).

### Within each user story

- Implement the handler before its test.
- Handlers across different commands are independent (different files) → parallel.

---

## Parallel Execution Examples

```bash
# Phase 2 — independent foundational files:
Task T002 (cli/charger.go)   Task T003 (cli/context.go)   Task T008 (cli/fake_charger_test.go)

# Phase 5 — all US3 handlers at once (different files):
T013 cmd_connect.go  T014 cmd_disconnect.go  T015 cmd_plugin.go  T016 cmd_unplug.go
T017 cmd_start.go    T018 cmd_stop.go        T019 cmd_meter.go

# Phase 5 — all US3 handler tests at once (after their handlers):
T020..T026  (cmd_*_test.go — one per command)

# Phase 6 + Phase 7 handlers can also all run in parallel once Phase 2 is done.
```

---

## Implementation Strategy

### MVP First

1. Phase 1 Setup → Phase 2 Foundational (the harness) → Phase 3 (US1 harness verification).
2. Add Phase 4 (US2 parsing) and Phase 5 (US3 lifecycle) — this delivers tested core OCPP commands.
3. **STOP and VALIDATE**: `go test ./cli/` green for the harness + lifecycle commands = a viable MVP proving the testing approach works end-to-end.

### Incremental Delivery

1. Foundation + US1 → harness proven.
2. + US2 → parsing locked.
3. + US3 → core lifecycle commands tested (MVP).
4. + US4 → parameter commands tested.
5. + US5 → info/meta commands tested; every command now covered.
6. + Polish → coverage guard + flakiness/perf/behavior checks.

---

## Notes

- [P] = different file, no dependency on an incomplete task.
- Every command verb has a dedicated handler file AND a dedicated test file (user directive).
- Assert on stable output substrings (per contracts/commands.md), not the entire byte stream.
- Refactor is behavior-preserving — output strings and state effects match the current `main.go`.
- Commit after each handler+test pair for clean, reviewable increments.
