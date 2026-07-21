package cli

import "fmt"

func init() { register("soc", handleSoc) }

// handleSoc sets the State of Charge from a numeric argument.
func handleSoc(ctx *CommandContext, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(ctx.Out, "Usage: soc <0-100>")
		return
	}
	var soc float64
	if _, err := fmt.Sscanf(args[0], "%f", &soc); err != nil {
		fmt.Fprintf(ctx.Out, "Error: invalid SOC value: %s\n", args[0])
		return
	}
	if err := ctx.Charger.SetSOC(soc); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
	} else {
		fmt.Fprintf(ctx.Out, "SOC set to: %.1f%%\n", soc)
	}
}
