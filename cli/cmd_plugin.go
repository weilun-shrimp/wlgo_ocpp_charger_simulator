package cli

import "fmt"

func init() { register("plugin", handlePlugin) }

// handlePlugin simulates a car plugging in.
func handlePlugin(ctx *CommandContext, args []string) {
	if err := ctx.Charger.Plugin(); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
	} else {
		fmt.Fprintln(ctx.Out, "Car plugged in (Preparing)")
	}
}
