# OCPP Charger Simulator

A command-line OCPP charger simulator written in Go. Supports OCPP 1.6 and 2.0.1 protocols.

## Quick Start

1. Copy the example configuration:
```bash
cp config.example.yaml config.yaml
```

2. Edit `config.yaml` with your server URL and settings:
```yaml
ocpp_version: "1.6"
charger_id: "CHARGER001"
server_url: "ws://localhost:8080/ocpp/CHARGER001"
voltage: 230
max_current: 32
max_power: 22000
```

3. Run the simulator:
```bash
go run main.go --config config.yaml
```

## Commands

| Command | Description |
|---------|-------------|
| `help` | Show available commands |
| `connect` | Connect to OCPP server |
| `disconnect` | Disconnect from server |
| `plugin` | Simulate car plug in (Available -> Preparing) |
| `unplug` | Simulate car unplug (-> Available) |
| `start <idTag>` | Start transaction (requires Preparing status) |
| `stop [reason]` | Stop transaction (reason: Local, Remote, etc.) |
| `status <status>` | Set charger status |
| `plate <plate>` | Send license plate via DataTransfer |
| `meter` | Send MeterValues manually |
| `soc <0-100>` | Set State of Charge |
| `current <amps>` | Set charging current (local control) |
| `info` | Show current charger status |

## Typical Charging Flow

```
> connect                    # Connect to OCPP server
> plugin                     # Car plugs in -> Preparing
> plate ABC-1234             # Send license plate (optional)
> start user123              # Start charging -> Charging (SOC auto-increases)
> stop                       # Stop charging -> Finishing
> unplug                     # Car unplugs -> Available
```

## Valid Statuses

**OCPP 1.6:** Available, Preparing, Charging, SuspendedEVSE, SuspendedEV, Finishing, Reserved, Unavailable, Faulted

**OCPP 2.0.1:** Available, Occupied, Reserved, Unavailable, Faulted

## Configuration

| Field | Description | Default |
|-------|-------------|---------|
| `ocpp_version` | "1.6" or "2.0.1" | Required |
| `charger_id` | Charger identity | Required |
| `server_url` | WebSocket URL (ws:// or wss://) | Required |
| `max_current` | Maximum current (A) | Required |
| `max_power` | Maximum power (W) | Required |
| `min_current` | Minimum current (A) | 0 |
| `min_power` | Minimum power (W) | 0 |
| `voltage` | Voltage (V) for power calculation | 230 |
| `connector_id` | Connector ID | 1 |
| `initial_status` | Initial charger status | Available |
| `initial_soc` | Initial State of Charge (%) | 20 |
| `battery_capacity` | Battery capacity (Wh) | 60000 |
| `meter_values_interval` | MeterValues interval (seconds) | 30 |

### TLS Configuration

For secure connections (wss://), add TLS config:

```yaml
tls:
  # Option 1: Trust CA certificate
  ca_file: "/path/to/ca.crt"

  # Option 2: Trust specific server cert (self-signed)
  server_cert_file: "/path/to/server.crt"

  # Option 3: Skip verification (insecure)
  skip_verify: true

  # Client certificate (mTLS)
  cert_file: "/path/to/client.crt"
  key_file: "/path/to/client.key"
```

## Features

- Supports OCPP 1.6 and 2.0.1
- Interactive CLI
- Current control (local via CLI, remote via SetChargingProfile)
- Auto SOC increase during charging
- License plate sending via DataTransfer
- TLS/mTLS support
- Offline operation (commands work without server connection)

## OCPP Messages Supported

| Message | Direction | Description |
|---------|-----------|-------------|
| BootNotification | CP -> CS | Sent on connect |
| StatusNotification | CP -> CS | Status changes |
| StartTransaction | CP -> CS | Start charging (1.6) |
| StopTransaction | CP -> CS | Stop charging (1.6) |
| TransactionEvent | CP -> CS | Transaction events (2.0.1) |
| MeterValues | CP -> CS | Energy/power readings |
| DataTransfer | CP -> CS | License plate, custom data |
| Heartbeat | CP -> CS | Keep-alive |
| RemoteStartTransaction | CS -> CP | Remote start (handled) |
| RemoteStopTransaction | CS -> CP | Remote stop (handled) |
| SetChargingProfile | CS -> CP | Remote current (1.6: A only, 2.0.1: A or W) |

## Build

```bash
go build -o charger-simulator main.go
./charger-simulator --config config.yaml
```
