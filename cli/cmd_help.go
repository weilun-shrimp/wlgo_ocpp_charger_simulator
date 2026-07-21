package cli

import (
	"fmt"
	"io"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/config"
)

func init() { register("help", handleHelp) }

// handleHelp prints the list of available commands and the valid-status list
// for the configured OCPP version.
func handleHelp(ctx *CommandContext, args []string) {
	printHelp(ctx.Out, ctx.Config)
}

func printHelp(out io.Writer, cfg *config.Config) {
	fmt.Fprintln(out, "Available commands:")
	fmt.Fprintln(out, "  help              - Show this help message")
	fmt.Fprintln(out, "  connect           - Connect to OCPP server")
	fmt.Fprintln(out, "  disconnect        - Disconnect from server")
	fmt.Fprintln(out, "  plugin            - Simulate car plug in (Preparing)")
	fmt.Fprintln(out, "  unplug            - Simulate car unplug (Available)")
	fmt.Fprintln(out, "  start <idTag>     - Start a transaction (requires Preparing status)")
	fmt.Fprintln(out, "  stop [reason]     - Stop the current transaction (reason: Local, Remote, etc.)")
	fmt.Fprintln(out, "  status <status>   - Set charger status (type 'status' for valid values)")
	fmt.Fprintln(out, "  plate <plate>     - Send license plate via DataTransfer")
	fmt.Fprintln(out, "  meter             - Send MeterValues")
	fmt.Fprintln(out, "  soc <0-100>       - Set State of Charge")
	fmt.Fprintf(out, "  current <amps>    - Set charging current (0-%.1f A, 0 = SuspendedEVSE)\n", cfg.MaxCurrent)
	fmt.Fprintf(out, "  power <watts>     - Set charging power (0-%.1f W, 0 = SuspendedEVSE)\n", cfg.MaxPower)
	fmt.Fprintln(out, "  info              - Show current charger status")
	fmt.Fprintln(out, "  quit/exit         - Exit the simulator (use Ctrl+C)")
	fmt.Fprintln(out)
	if cfg.IsOCPP16() {
		fmt.Fprintln(out, "Valid statuses (OCPP 1.6): Available, Preparing, Charging, SuspendedEVSE, SuspendedEV, Finishing, Reserved, Unavailable, Faulted")
	} else {
		fmt.Fprintln(out, "Valid statuses (OCPP 2.0.1): Available, Occupied, Reserved, Unavailable, Faulted")
	}
}
