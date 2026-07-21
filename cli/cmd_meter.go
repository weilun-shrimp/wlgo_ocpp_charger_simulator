package cli

import "fmt"

func init() { register("meter", handleMeter) }

// handleMeter sends a MeterValues message.
func handleMeter(ctx *CommandContext, args []string) {
	if err := ctx.Charger.MeterValues(); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
	} else {
		fmt.Fprintln(ctx.Out, "MeterValues updated")
	}
}
