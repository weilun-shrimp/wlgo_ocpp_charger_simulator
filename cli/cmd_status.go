package cli

import "fmt"

func init() { register("status", handleStatus) }

// handleStatus sets the charger status, or prints usage plus the valid-status
// list for the configured OCPP version when no argument is given.
func handleStatus(ctx *CommandContext, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(ctx.Out, "Usage: status <status>")
		if ctx.Config.IsOCPP16() {
			fmt.Fprintln(ctx.Out, "Valid statuses (OCPP 1.6): Available, Preparing, Charging, SuspendedEVSE, SuspendedEV, Finishing, Reserved, Unavailable, Faulted")
		} else {
			fmt.Fprintln(ctx.Out, "Valid statuses (OCPP 2.0.1): Available, Occupied, Reserved, Unavailable, Faulted")
		}
		return
	}
	status := args[0]
	if err := ctx.Charger.SetStatus(status); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
	} else {
		fmt.Fprintf(ctx.Out, "Status updated to: %s\n", status)
	}
}
