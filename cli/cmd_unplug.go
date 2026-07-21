package cli

import "fmt"

func init() { register("unplug", handleUnplug) }

// handleUnplug simulates a car unplugging.
func handleUnplug(ctx *CommandContext, args []string) {
	if err := ctx.Charger.Unplug(); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
	} else {
		fmt.Fprintln(ctx.Out, "Car unplugged (Available)")
	}
}
