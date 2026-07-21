package cli

import "fmt"

func init() { register("current", handleCurrent) }

// handleCurrent sets the charging current, or reports usage plus the present
// value when no argument is given.
func handleCurrent(ctx *CommandContext, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(ctx.Out, "Usage: current <amperes> (0-%.1f A, 0 = SuspendedEVSE)\n", ctx.Config.MaxCurrent)
		fmt.Fprintf(ctx.Out, "Current: %.1f A\n", ctx.Charger.GetCurrent())
		return
	}
	var current float64
	if _, err := fmt.Sscanf(args[0], "%f", &current); err != nil {
		fmt.Fprintf(ctx.Out, "Error: invalid current value: %s\n", args[0])
		return
	}
	if err := ctx.Charger.SetCurrent(current); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
	} else {
		fmt.Fprintf(ctx.Out, "Current set to: %.1f A\n", current)
	}
}
