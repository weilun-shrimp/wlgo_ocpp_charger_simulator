# Feature Specification: Interactive Command Test Coverage

**Feature Branch**: `001-interactive-command-tests`

**Created**: 2026-07-21

**Status**: Draft

**Input**: User description: "There is no any testing in main.go interactiveLoop comand line behavior. I want each one of them all have it own testing."

## User Scenarios & Testing *(mandatory)*

The "users" of this feature are the developers, testers, and QA engineers who maintain
the simulator. The value delivered is a trustworthy automated test suite that pins the
behavior of every interactive command so that regressions are caught before release.

Each command handled by the interactive command loop is a separately verifiable unit.
The stories below group those commands so each group can be delivered and validated on
its own, but the acceptance bar is the same throughout: **every interactive command has
at least one dedicated, deterministic automated test that asserts its observable output
and its effect on charger state.**

### User Story 1 - Testable command dispatch harness (Priority: P1)

As a maintainer, I can drive any interactive command with a controlled line of input and
capture everything the command prints, without touching the real keyboard, the real
network, or an external OCPP server, and without the program blocking on an infinite
input loop.

**Why this priority**: Nothing else in this feature is possible until a single command
can be executed in isolation and its output observed. This is the foundation every
per-command test builds on, so it is the minimum viable slice.

**Independent Test**: Feed a single known command line into the dispatch harness, capture
the produced output, and assert the captured text matches the expected message — proving
commands can be exercised and observed in a test.

**Acceptance Scenarios**:

1. **Given** a simulator instance in a known state, **When** a test submits one command
   line to the dispatcher, **Then** the dispatcher processes exactly that command and the
   test can read back the full text it printed.
2. **Given** the dispatcher processing test input, **When** the input is exhausted,
   **Then** processing returns control to the test instead of blocking forever.
3. **Given** a command that would otherwise require a live server connection, **When** it
   is run in a test, **Then** it runs against a test double and never contacts an external
   network endpoint.

---

### User Story 2 - Input parsing behavior is verified (Priority: P1)

As a maintainer, the shared parsing rules that sit in front of every command — trimming
whitespace, ignoring empty lines, lower-casing the command word, splitting arguments, and
rejecting unknown commands — are covered by their own tests.

**Why this priority**: These rules apply to every command. A defect here breaks all
commands at once, so they must be locked down early and independently.

**Independent Test**: Submit blank lines, whitespace-only lines, mixed-case command
words, and a nonsense command, and assert each is handled as specified without invoking
any command action.

**Acceptance Scenarios**:

1. **Given** an empty line or a whitespace-only line, **When** it is submitted, **Then**
   no command runs and no error is produced.
2. **Given** a command word in mixed or upper case (e.g. `CONNECT`, `Info`), **When** it
   is submitted, **Then** it is treated identically to its lower-case form.
3. **Given** an unrecognized command word, **When** it is submitted, **Then** the output
   states the command is unknown and points the user to `help`.
4. **Given** a command followed by extra arguments, **When** it is submitted, **Then** the
   first token is used as the command and the remaining tokens are available as arguments.

---

### User Story 3 - Session lifecycle commands are verified (Priority: P1)

As a maintainer, the commands that drive the connection and charging session lifecycle —
`connect`, `disconnect`, `plugin`, `unplug`, `start`, `stop`, `meter` — each have their
own tests covering both their success path and their guarded/error path.

**Why this priority**: These commands exercise the core OCPP flows the simulator exists to
reproduce; per the constitution, core protocol flows must be covered by automated tests.

**Independent Test**: Run each lifecycle command from a known state and assert both the
printed message and the resulting charger state; then run it from a state where it should
be refused and assert the guard message.

**Acceptance Scenarios**:

1. **Given** a disconnected simulator, **When** `connect` succeeds, **Then** the output
   confirms the connection and the follow-up boot/status steps are attempted.
2. **Given** an already-connected simulator, **When** `connect` is submitted, **Then** the
   output reports it is already connected and no second connection is attempted.
3. **Given** a disconnected simulator, **When** `disconnect` is submitted, **Then** the
   output reports it is not connected.
4. **Given** a connected simulator, **When** `disconnect` is submitted, **Then** the
   output confirms disconnection.
5. **Given** any simulator, **When** `plugin` / `unplug` / `meter` succeed, **Then** each
   prints its confirmation message; **When** the underlying action fails, **Then** each
   prints the returned error.
6. **Given** `start` with no id tag, **When** submitted, **Then** the output shows the
   usage hint; **When** given an id tag, **Then** it reports the transaction started or
   prints the returned error.
7. **Given** `stop` with no reason, **When** submitted, **Then** it uses the default reason
   and reports the transaction stopped; **When** given a reason, **Then** that reason is
   used.

---

### User Story 4 - State and parameter commands are verified (Priority: P2)

As a maintainer, the commands that read or change charger parameters — `status`, `soc`,
`current`, `power`, `plate` — each have their own tests covering the usage/help path, the
invalid-input path, and the successful-update path.

**Why this priority**: These commands accept user arguments and perform validation, so they
carry the highest risk of parsing and boundary defects; they come after the lifecycle core.

**Independent Test**: For each command, submit it with no argument, with a malformed
argument, and with a valid argument, asserting the usage text, the specific error message,
and the confirmation message respectively.

**Acceptance Scenarios**:

1. **Given** `status` / `soc` / `current` / `power` / `plate` with no argument, **When**
   submitted, **Then** the output shows that command's usage hint (and current value where
   the command reports one).
2. **Given** `soc` / `current` / `power` with a non-numeric argument, **When** submitted,
   **Then** the output reports the value is invalid and no state change occurs.
3. **Given** a valid argument, **When** submitted, **Then** the command reports the updated
   value and the charger state reflects the change (or prints the returned error if the
   action is refused).
4. **Given** the configured OCPP version, **When** `status` shows its valid-status list,
   **Then** the list matches the set valid for that version (1.6 vs 2.0.1).

---

### User Story 5 - Informational and meta commands are verified (Priority: P2)

As a maintainer, the read-only and meta commands — `help`, `info`, `quit`/`exit` — each
have their own tests.

**Why this priority**: These do not change protocol state, so they are lower risk, but they
are part of "each command has its own test" and are cheap to cover.

**Independent Test**: Run each command and assert its output: `help` lists the commands and
the version-appropriate status set; `info` reports the current charger fields; `quit`/`exit`
prints the Ctrl+C guidance.

**Acceptance Scenarios**:

1. **Given** any state, **When** `help` is submitted, **Then** the output lists every
   available command and the valid-status list for the configured OCPP version.
2. **Given** a known charger state, **When** `info` is submitted, **Then** the output
   reports connection, status, charging flag, voltage, current, power, and SOC, and
   includes the license-plate line only when a plate is set.
3. **Given** any state, **When** `quit` or `exit` is submitted, **Then** the output tells
   the user to use Ctrl+C and the program continues running.

---

### Edge Cases

- A command whose underlying charger action returns an error MUST surface that error text
  rather than a success message.
- Numeric commands (`soc`, `current`, `power`) given empty, non-numeric, or partially
  numeric input MUST be rejected with the invalid-value message and cause no state change.
- The `info` license-plate line MUST appear only when a plate is set and be absent
  otherwise.
- Commands run against a test double MUST behave deterministically regardless of whether a
  real server is reachable.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The interactive command behavior MUST be exercisable by automated tests
  that supply command input and capture command output without reading real stdin, opening
  a real network connection, or requiring an external OCPP server.
- **FR-002**: Command processing MUST terminate when its supplied input is exhausted so a
  test can run a finite sequence of commands and then assert results.
- **FR-003**: Every command accepted by the interactive loop MUST have at least one
  dedicated automated test. The commands are: `help`, `connect`, `disconnect`, `plugin`,
  `unplug`, `status`, `start`, `stop`, `meter`, `plate`, `soc`, `current`, `power`,
  `info`, `quit`, and `exit`.
- **FR-004**: The unknown-command path and the empty/whitespace-only input path MUST each
  have their own tests.
- **FR-005**: Shared input parsing — whitespace trimming, empty-line skipping,
  case-insensitive command matching, and argument splitting — MUST be covered by tests.
- **FR-006**: For each command that has both a success path and a guarded/error path (e.g.
  `connect` when already connected, `disconnect` when not connected, missing-argument
  usage hints, invalid numeric input), tests MUST cover both paths.
- **FR-007**: Tests MUST assert observable behavior: the exact or substring-matched text a
  command prints, and, where the command changes charger state, the resulting state.
- **FR-008**: Tests for version-dependent output (`help` and the `status` usage list) MUST
  verify the correct set is shown for both OCPP 1.6 and OCPP 2.0.1 configurations.
- **FR-009**: The test suite MUST be deterministic — repeated runs with the same inputs
  produce the same results — and MUST run via the project's standard test command with no
  manual setup or external services.
- **FR-010**: When a command's underlying action fails, the corresponding test MUST assert
  that the failure is surfaced to the user rather than reported as success.
- **FR-011**: Adding a new interactive command in the future MUST require adding a
  corresponding test for it; the suite's structure MUST make an untested command
  visible (e.g. a coverage list or table that fails when a command is missing).

### Key Entities

- **Interactive Command**: A single verb accepted at the prompt (e.g. `connect`, `soc`),
  with optional arguments, a printed response, and an optional effect on charger state.
- **Command Test Case**: A named scenario pairing an input line (and starting charger
  state) with the expected printed output and expected resulting state.
- **Charger Test Double**: A stand-in for the charger's server-facing behavior that lets
  connection- and protocol-dependent commands run deterministically without a live server.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of the interactive commands listed in FR-003 have at least one dedicated
  automated test.
- **SC-002**: Both the success path and the guarded/error path are tested for every command
  that has two paths (as enumerated in FR-006).
- **SC-003**: The full interactive-command test suite runs to completion in under 10
  seconds on a developer machine and requires no network access or external server.
- **SC-004**: Running the suite twice in a row yields identical pass/fail results (zero
  flakiness) across at least 10 consecutive runs.
- **SC-005**: A deliberately introduced regression in any single command's output or state
  effect causes at least one test to fail, demonstrating the suite's protective value.
- **SC-006**: A maintainer can determine which commands are covered by reading the test
  suite, and an interactive command with no test is detectable from the suite itself.

## Assumptions

- Tests are written with the project's existing language-standard test tooling (`go test`)
  and live alongside the code they cover; no new external test framework is introduced.
- "Its own testing" means each command has a dedicated, independently runnable test case
  (table-driven cases count as dedicated per-command cases); it does not require a separate
  test binary per command.
- To satisfy the "no live server / deterministic" requirement, the command layer will be
  made observable and its server-facing dependency substitutable via a test double; the
  exact seam (interface extraction, injected reader/writer, or refactor of the dispatch out
  of the blocking loop) is an implementation decision left to planning.
- Assertions may match on stable substrings of command output rather than the entire byte
  stream, so cosmetic wording tweaks do not force brittle test churn, while still pinning
  the meaningful content.
- Commands whose success depends on prior charger state (e.g. `start` requires `Preparing`)
  are tested by first placing the charger test double in the required state.
