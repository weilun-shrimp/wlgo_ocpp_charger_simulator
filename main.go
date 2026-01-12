package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/charger"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("OCPP Charger Simulator")
	log.Printf("======================")
	log.Printf("Charger ID: %s", cfg.ChargerID)
	log.Printf("OCPP Version: %s", cfg.OCPPVersion)
	log.Printf("Server URL: %s", cfg.ServerURL)
	log.Printf("Voltage: %.1f V", cfg.Voltage)
	log.Printf("Max Current: %.1f A", cfg.MaxCurrent)
	log.Printf("Max Power: %.1f W", cfg.MaxPower)
	log.Printf("Initial Status: %s", cfg.InitialStatus)
	log.Printf("Initial SOC: %.1f%%", cfg.InitialSOC)
	log.Printf("Battery Capacity: %.0f Wh", cfg.BatteryCapacity)

	// Create charger
	sim, err := charger.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create charger: %v", err)
	}
	defer sim.Close()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start interactive command loop
	go interactiveLoop(sim, cfg)

	log.Println("Charger simulator ready. Type 'connect' to connect to server, 'help' for commands.")

	// Wait for shutdown signal
	<-sigCh
	log.Println("Shutting down...")
}

func interactiveLoop(sim *charger.Charger, cfg *config.Config) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		input = strings.TrimSpace(input)
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmd := strings.ToLower(parts[0])

		switch cmd {
		case "help":
			printHelp(cfg)

		case "connect":
			if sim.IsConnected() {
				fmt.Println("Already connected")
				continue
			}
			if err := sim.Connect(); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("Connected to server")

			// Send BootNotification
			if err := sim.BootNotification(); err != nil {
				fmt.Printf("BootNotification failed: %v\n", err)
				continue
			}

			// Send initial StatusNotification
			if err := sim.StatusNotification(cfg.InitialStatus); err != nil {
				fmt.Printf("StatusNotification failed: %v\n", err)
			}

		case "disconnect":
			if !sim.IsConnected() {
				fmt.Println("Not connected")
				continue
			}
			sim.Disconnect()
			fmt.Println("Disconnected from server")

		case "plugin":
			if err := sim.Plugin(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Car plugged in (Preparing)")
			}

		case "unplug":
			if err := sim.Unplug(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Car unplugged (Available)")
			}

		case "status":
			if len(parts) < 2 {
				if cfg.IsOCPP16() {
					fmt.Println("Usage: status <status>")
					fmt.Println("Valid statuses (OCPP 1.6): Available, Preparing, Charging, SuspendedEVSE, SuspendedEV, Finishing, Reserved, Unavailable, Faulted")
				} else {
					fmt.Println("Usage: status <status>")
					fmt.Println("Valid statuses (OCPP 2.0.1): Available, Occupied, Reserved, Unavailable, Faulted")
				}
				continue
			}
			status := parts[1]
			if err := sim.SetStatus(status); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Status updated to: %s\n", status)
			}

		case "start":
			if len(parts) < 2 {
				fmt.Println("Usage: start <idTag>")
				continue
			}
			idTag := parts[1]
			if err := sim.StartTransaction(idTag); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Transaction started")
			}

		case "stop":
			reason := "Local"
			if len(parts) >= 2 {
				reason = parts[1]
			}
			if err := sim.StopTransaction(reason); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Transaction stopped")
			}

		case "meter":
			if err := sim.MeterValues(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("MeterValues updated")
			}

		case "plate":
			if len(parts) < 2 {
				fmt.Println("Usage: plate <license_plate>")
				continue
			}
			plate := parts[1]
			if err := sim.SetLicensePlateAndSend(plate); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("License plate set: %s\n", plate)
			}

		case "soc":
			if len(parts) < 2 {
				fmt.Println("Usage: soc <0-100>")
				continue
			}
			var soc float64
			if _, err := fmt.Sscanf(parts[1], "%f", &soc); err != nil {
				fmt.Printf("Error: invalid SOC value: %s\n", parts[1])
				continue
			}
			if err := sim.SetSOC(soc); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("SOC set to: %.1f%%\n", soc)
			}

		case "current":
			if len(parts) < 2 {
				fmt.Printf("Usage: current <amperes> (min: %.1f A, max: %.1f A)\n", cfg.MinCurrent, cfg.MaxCurrent)
				fmt.Printf("Current: %.1f A\n", sim.GetCurrent())
				continue
			}
			var current float64
			if _, err := fmt.Sscanf(parts[1], "%f", &current); err != nil {
				fmt.Printf("Error: invalid current value: %s\n", parts[1])
				continue
			}
			if err := sim.SetCurrent(current); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Current set to: %.1f A\n", current)
			}

		case "info":
			fmt.Printf("Connected: %v\n", sim.IsConnected())
			fmt.Printf("Status: %s\n", sim.GetStatus())
			fmt.Printf("Charging: %v\n", sim.IsCharging())
			current := sim.GetCurrent()
			power := current * cfg.Voltage
			if power > cfg.MaxPower {
				power = cfg.MaxPower
			}
			fmt.Printf("Voltage: %.1f V\n", cfg.Voltage)
			fmt.Printf("Current: %.1f A\n", current)
			fmt.Printf("Power: %.1f W\n", power)
			fmt.Printf("SOC: %.1f%%\n", sim.GetSOC())
			if plate := sim.GetLicensePlate(); plate != "" {
				fmt.Printf("License Plate: %s\n", plate)
			}

		case "quit", "exit":
			fmt.Println("Use Ctrl+C to exit")

		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", cmd)
		}
	}
}

func printHelp(cfg *config.Config) {
	fmt.Println("Available commands:")
	fmt.Println("  help              - Show this help message")
	fmt.Println("  connect           - Connect to OCPP server")
	fmt.Println("  disconnect        - Disconnect from server")
	fmt.Println("  plugin            - Simulate car plug in (Preparing)")
	fmt.Println("  unplug            - Simulate car unplug (Available)")
	fmt.Println("  start <idTag>     - Start a transaction (requires Preparing status)")
	fmt.Println("  stop [reason]     - Stop the current transaction (reason: Local, Remote, etc.)")
	fmt.Println("  status <status>   - Set charger status (type 'status' for valid values)")
	fmt.Println("  plate <plate>     - Send license plate via DataTransfer")
	fmt.Println("  meter             - Send MeterValues")
	fmt.Println("  soc <0-100>       - Set State of Charge")
	fmt.Printf("  current <amps>    - Set charging current (%.1f-%.1f A)\n", cfg.MinCurrent, cfg.MaxCurrent)
	fmt.Println("  info              - Show current charger status")
	fmt.Println("  quit/exit         - Exit the simulator (use Ctrl+C)")
	fmt.Println()
	if cfg.IsOCPP16() {
		fmt.Println("Valid statuses (OCPP 1.6): Available, Preparing, Charging, SuspendedEVSE, SuspendedEV, Finishing, Reserved, Unavailable, Faulted")
	} else {
		fmt.Println("Valid statuses (OCPP 2.0.1): Available, Occupied, Reserved, Unavailable, Faulted")
	}
}
