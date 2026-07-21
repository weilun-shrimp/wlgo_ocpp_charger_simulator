package cli

import "fmt"

func init() { register("info", handleInfo) }

// handleInfo prints the current charger state. The license-plate line is shown
// only when a plate is set.
func handleInfo(ctx *CommandContext, args []string) {
	fmt.Fprintf(ctx.Out, "Connected: %v\n", ctx.Charger.IsConnected())
	fmt.Fprintf(ctx.Out, "Status: %s\n", ctx.Charger.GetStatus())
	fmt.Fprintf(ctx.Out, "Charging: %v\n", ctx.Charger.IsCharging())
	fmt.Fprintf(ctx.Out, "Voltage: %.1f V\n", ctx.Config.Voltage)
	fmt.Fprintf(ctx.Out, "Current: %.1f A\n", ctx.Charger.GetCurrent())
	fmt.Fprintf(ctx.Out, "Power: %.1f W\n", ctx.Charger.GetPower())
	fmt.Fprintf(ctx.Out, "SOC: %.1f%%\n", ctx.Charger.GetSOC())
	if plate := ctx.Charger.GetLicensePlate(); plate != "" {
		fmt.Fprintf(ctx.Out, "License Plate: %s\n", plate)
	}
}
