package cli

import "fmt"

func init() { register("power", handlePower) }

// handlePower sets the charging power, or reports usage plus the present value
// when no argument is given.
func handlePower(ctx *CommandContext, args []string) {
	if len(args) < 1 {
		fmt.Fprintf(ctx.Out, "Usage: power <watts> (0-%.1f W, 0 = SuspendedEVSE)\n", ctx.Config.MaxPower)
		fmt.Fprintf(ctx.Out, "Power: %.1f W\n", ctx.Charger.GetPower())
		return
	}
	var power float64
	if _, err := fmt.Sscanf(args[0], "%f", &power); err != nil {
		fmt.Fprintf(ctx.Out, "Error: invalid power value: %s\n", args[0])
		return
	}
	if err := ctx.Charger.SetPower(power); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
	} else {
		fmt.Fprintf(ctx.Out, "Power set to: %.1f W\n", power)
	}
}
