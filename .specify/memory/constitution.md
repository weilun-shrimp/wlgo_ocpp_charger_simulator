<!--
Sync Impact Report
==================
Version change: (template, unversioned) → 1.0.0
Bump rationale: Initial ratification of a concrete constitution replacing the
  unfilled template. First versioned adoption ⇒ MAJOR baseline 1.0.0.

Modified principles: N/A (initial adoption)
Added sections:
  - Purpose
  - Core Principles (9): Protocol Fidelity, CLI-First Design, Predictable
    Behavior, Configuration over Code, Clear Error Messages, CLI Output,
    Testability, Extensibility, Maintainability
  - Decision Order
  - Governance
Removed sections: All template placeholder scaffolding

Templates requiring updates:
  - .specify/templates/plan-template.md ............ ✅ aligned (generic
    Constitution Check gate; no principle names hardcoded)
  - .specify/templates/spec-template.md ............ ✅ aligned (no changes needed)
  - .specify/templates/tasks-template.md ........... ✅ aligned (tests remain
    OPTIONAL in template; Testability principle governs when they are required)
  - .claude/skills/speckit-*/ ...................... ✅ aligned (generic
    /speckit-* naming; no agent-specific references to fix)

Follow-up TODOs: None. Ratification date set to first adoption (2026-07-21).
-->

# OCPP Charger Simulator Constitution

## Purpose

This project builds a terminal-based OCPP charger simulator for developers, testers,
and QA engineers. The simulator MUST be predictable, scriptable, and faithful to the
OCPP specification while providing an excellent command-line experience.

## Core Principles

### I. Protocol Fidelity

- Correctness MUST be prioritized over convenience.
- OCPP behavior MUST be implemented according to the supported specification
  (OCPP 1.6 and OCPP 2.0.1).
- Any intentional deviation from the specification MUST be documented.
- Invalid protocol behavior MUST NEVER be the default.

Rationale: The simulator's value depends entirely on standing in for a real charger.
A subtly non-conformant simulator produces false confidence and wastes the time it was
built to save.

### II. CLI-First Design

- Every feature MUST be accessible from the command line.
- Commands MUST have clear names and consistent syntax.
- Output MUST be human-readable by default and machine-readable when requested
  (e.g., JSON).
- Commands MUST be composable and suitable for scripting.

Rationale: The primary users are developers, testers, and QA engineers who automate.
A composable CLI is what makes the tool usable inside test harnesses and pipelines.

### III. Predictable Behavior

- Running the same command with the same inputs MUST produce the same results unless
  randomness is explicitly enabled.
- Default values MUST be sensible and documented.
- State changes MUST be explicit.

Rationale: Determinism is the foundation of testability. Hidden state or unrequested
randomness makes failures impossible to reproduce and diagnose.

### IV. Configuration over Code

- Charger behavior MUST be configurable rather than hardcoded.
- Configuration files and environment variables MUST be supported where appropriate.
- Command-line arguments MUST override configuration values.

Rationale: Users need to model many charger variants and scenarios without editing or
recompiling source. A clear precedence order keeps overrides unambiguous.

### V. Clear Error Messages

Every error MUST explain:

- What failed
- Why it failed
- How to fix it

Rationale: A simulator is a diagnostic tool. An error that does not point toward a fix
shifts the debugging burden back onto the user it was meant to help.

### VI. CLI Output

- All CLI output MUST be human-readable.
- Success messages MUST clearly indicate completed actions.
- Error messages MUST explain what failed and, where possible, SHOULD suggest how to
  resolve the issue.
- Output MUST remain consistent across commands.

Rationale: Consistent, legible output lets users trust what they see and lets scripts
parse it reliably. Consistency across commands lowers the cost of learning the tool.

### VII. Testability

- Every command MUST be testable.
- Core protocol flows MUST be covered by automated tests.
- Regression tests MUST accompany bug fixes.

Rationale: Protocol correctness cannot be maintained by inspection alone. Automated
coverage of core flows is what keeps Protocol Fidelity true over time.

### VIII. Extensibility

The architecture MUST make it straightforward to add:

- New OCPP versions
- Additional commands
- New simulation scenarios
- Vendor-specific extensions
- Multiple simulated chargers

Rationale: OCPP evolves and real deployments carry vendor quirks. Designing for
extension avoids rewrites each time a new version or scenario is required.

### IX. Maintainability

- Business logic MUST be kept separate from CLI parsing.
- Duplicated protocol logic MUST be avoided.
- Small, focused modules MUST be favored.

Rationale: Separating protocol logic from the interface layer keeps both independently
testable and lets either evolve without destabilizing the other.

## Decision Order

When trade-offs arise, priorities MUST be weighed in this order:

1. OCPP correctness
2. User experience
3. Simplicity
4. Extensibility
5. Performance

A lower-ranked concern MUST NOT be pursued at the expense of a higher-ranked one
without explicit, documented justification.

## Governance

- This constitution supersedes other development practices when they conflict.
- Amendments MUST be proposed via pull request, documenting the change, its rationale,
  and any migration impact, and MUST be approved before merge.
- Versioning follows semantic versioning:
  - MAJOR: backward-incompatible governance or principle removals/redefinitions.
  - MINOR: a new principle or section, or materially expanded guidance.
  - PATCH: clarifications, wording, or non-semantic refinements.
- All pull requests and reviews MUST verify compliance with these principles.
  Deviations MUST be justified in writing or the change MUST be revised to comply.
- Complexity that violates a principle MUST be justified against the Decision Order or
  removed.

**Version**: 1.0.0 | **Ratified**: 2026-07-21 | **Last Amended**: 2026-07-21
