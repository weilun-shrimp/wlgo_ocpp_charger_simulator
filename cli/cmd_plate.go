package cli

import "fmt"

func init() { register("plate", handlePlate) }

// handlePlate sends the EV license plate via DataTransfer.
func handlePlate(ctx *CommandContext, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(ctx.Out, "Usage: plate <license_plate>")
		return
	}
	plate := args[0]
	if err := ctx.Charger.SetLicensePlateAndSend(plate); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
	} else {
		fmt.Fprintf(ctx.Out, "License plate set: %s\n", plate)
	}
}
