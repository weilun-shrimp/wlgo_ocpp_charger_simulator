package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/charger"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/cli"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/config"
)

// Compile-time assertion that the concrete charger satisfies the CLI's Charger
// interface. Keeps the interface honest without the cli package importing charger.
var _ cli.Charger = (*charger.Charger)(nil)

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
	go cli.Run(sim, cfg, os.Stdin, os.Stdout)

	log.Println("Charger simulator ready. Type 'connect' to connect to server, 'help' for commands.")

	// Wait for shutdown signal
	<-sigCh
	log.Println("Shutting down...")
}
