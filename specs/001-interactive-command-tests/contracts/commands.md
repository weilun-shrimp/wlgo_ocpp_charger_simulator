# CLI Command Contract: Interactive Commands

This contract enumerates every interactive command, its arguments, its charger dependency, and
its observable output — the behavior that must be preserved by the refactor and pinned by tests.
Output strings are the current wording (Principle V/VI: preserved verbatim). Tests SHOULD assert
on the stable substrings shown in **Output** rather than the full byte stream.

Legend — **Guard**: precondition path that prints a specific message instead of acting.

## Parsing (applies to all commands)

| Rule | Input | Behavior |
|------|-------|----------|
| Empty/whitespace | `""`, `"   "` | No command runs; no output beyond the prompt |
| Case-insensitive verb | `CONNECT`, `Info` | Treated as `connect`, `info` |
| Argument split | `soc 55` | verb=`soc`, args=`["55"]` |
| Unknown verb | `frobnicate` | `Unknown command: frobnicate. Type 'help' for available commands.` |

## help

- **Args**: none
- **Charger**: none (reads `Config.OCPPVersion`)
- **Output**: full command list; then the valid-status line for the configured version:
  - 1.6: `Valid statuses (OCPP 1.6): Available, Preparing, Charging, SuspendedEVSE, SuspendedEV, Finishing, Reserved, Unavailable, Faulted`
  - 2.0.1: `Valid statuses (OCPP 2.0.1): Available, Occupied, Reserved, Unavailable, Faulted`
- **Test**: `TestHandleHelp` — both versions.

## connect

- **Args**: none
- **Guard**: if `IsConnected()` → `Already connected` (no `Connect`).
- **Success**: `Connect()` ok → `Connected to server`, then `BootNotification()` then
  `StatusNotification(GetStatus())` attempted.
- **Errors**: `Connect` err → `Error: <err>`; `BootNotification` err → `BootNotification failed: <err>`;
  `StatusNotification` err → `StatusNotification failed: <err>`.
- **Test**: `TestHandleConnect` — already-connected guard, success chain, connect error.

## disconnect

- **Args**: none
- **Guard**: if not connected → `Not connected`.
- **Success**: `Disconnect()` → `Disconnected from server`.
- **Test**: `TestHandleDisconnect` — both paths.

## plugin

- **Args**: none
- **Success**: `Plugin()` ok → `Car plugged in (Preparing)`.
- **Error**: `Plugin()` err → `Error: <err>`.
- **Test**: `TestHandlePlugin` — success + error.

## unplug

- **Args**: none
- **Success**: `Unplug()` ok → `Car unplugged (Available)`.
- **Error**: `Unplug()` err → `Error: <err>`.
- **Test**: `TestHandleUnplug` — success + error.

## status

- **Args**: `<status>` (optional)
- **Guard/usage**: no arg → `Usage: status <status>` + version-appropriate valid-status line.
- **Success**: `SetStatus(arg)` ok → `Status updated to: <arg>`.
- **Error**: `SetStatus` err → `Error: <err>`.
- **Test**: `TestHandleStatus` — usage (both versions), success, error.

## start

- **Args**: `<idTag>` (required)
- **Usage**: no arg → `Usage: start <idTag>`.
- **Success**: `StartTransaction(idTag)` ok → `Transaction started`.
- **Error**: err → `Error: <err>`.
- **Test**: `TestHandleStart` — usage, success, error.

## stop

- **Args**: `[reason]` (optional; default `Local`)
- **Success**: `StopTransaction(reason)` ok → `Transaction stopped` (reason defaults to `Local`).
- **Error**: err → `Error: <err>`.
- **Test**: `TestHandleStop` — default reason, explicit reason, error.

## meter

- **Args**: none
- **Success**: `MeterValues()` ok → `MeterValues updated`.
- **Error**: err → `Error: <err>`.
- **Test**: `TestHandleMeter` — success + error.

## plate

- **Args**: `<license_plate>` (required)
- **Usage**: no arg → `Usage: plate <license_plate>`.
- **Success**: `SetLicensePlateAndSend(arg)` ok → `License plate set: <arg>`.
- **Error**: err → `Error: <err>`.
- **Test**: `TestHandlePlate` — usage, success, error.

## soc

- **Args**: `<0-100>` (required, numeric)
- **Usage**: no arg → `Usage: soc <0-100>`.
- **Invalid**: non-numeric → `Error: invalid SOC value: <arg>` (no state change).
- **Success**: `SetSOC(v)` ok → `SOC set to: <v>%` (one decimal).
- **Error**: `SetSOC` err (out of range) → `Error: <err>`.
- **Test**: `TestHandleSoc` — usage, invalid, success, range error.

## current

- **Args**: `<amperes>` (optional)
- **Usage**: no arg → `Usage: current <amperes> (0-<MaxCurrent> A, 0 = SuspendedEVSE)` and
  `Current: <GetCurrent()> A`.
- **Invalid**: non-numeric → `Error: invalid current value: <arg>`.
- **Success**: `SetCurrent(v)` ok → `Current set to: <v> A`.
- **Error**: err → `Error: <err>`.
- **Test**: `TestHandleCurrent` — usage, invalid, success, error.

## power

- **Args**: `<watts>` (optional)
- **Usage**: no arg → `Usage: power <watts> (0-<MaxPower> W, 0 = SuspendedEVSE)` and
  `Power: <GetPower()> W`.
- **Invalid**: non-numeric → `Error: invalid power value: <arg>`.
- **Success**: `SetPower(v)` ok → `Power set to: <v> W`.
- **Error**: err → `Error: <err>`.
- **Test**: `TestHandlePower` — usage, invalid, success, error.

## info

- **Args**: none
- **Output**: `Connected`, `Status`, `Charging`, `Voltage`, `Current`, `Power`, `SOC` lines from
  the charger/config; `License Plate: <plate>` line **only** when `GetLicensePlate() != ""`.
- **Test**: `TestHandleInfo` — plate present vs absent.

## quit / exit

- **Args**: none
- **Output**: `Use Ctrl+C to exit` (loop keeps running).
- **Test**: `TestHandleQuit` — both aliases.

## Coverage guard

- **Test**: `TestRegistryCoversAllCommands` — asserts the registry has a handler for each
  canonical verb (help, connect, disconnect, plugin, unplug, status, start, stop, meter, plate,
  soc, current, power, info, quit, exit). Fails if a command is added without registration.
