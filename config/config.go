package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// TLSConfig holds TLS certificate configuration
type TLSConfig struct {
	CAFile         string `yaml:"ca_file"`          // CA certificate to verify server cert chain
	ServerCertFile string `yaml:"server_cert_file"` // Trusted server certificate (for self-signed certs)
	CertFile       string `yaml:"cert_file"`        // Client certificate
	KeyFile        string `yaml:"key_file"`         // Client private key
	SkipVerify     bool   `yaml:"skip_verify"`      // Skip server certificate verification (insecure)
}

// Config holds the charger simulator configuration
type Config struct {
	OCPPVersion         string     `yaml:"ocpp_version"`
	ChargerID           string     `yaml:"charger_id"`
	ServerURL           string     `yaml:"server_url"`
	TLS                 *TLSConfig `yaml:"tls"`
	InitialStatus       string     `yaml:"initial_status"`
	MaxCurrent          float64    `yaml:"max_current"`
	MaxPower            float64    `yaml:"max_power"`
	MinCurrent          float64    `yaml:"min_current"`
	MinPower            float64    `yaml:"min_power"`
	Voltage             float64    `yaml:"voltage"` // Voltage in V (for power calculation)
	ConnectorID         int        `yaml:"connector_id"`
	MeterValuesInterval int        `yaml:"meter_values_interval"`
	// EV Battery simulation
	InitialSOC      float64 `yaml:"initial_soc"`      // Initial State of Charge (0-100%)
	BatteryCapacity float64 `yaml:"battery_capacity"` // Battery capacity in Wh
}

// Load reads and parses the configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{
		// Set defaults
		InitialStatus:       "Available",
		MinCurrent:          0,
		MinPower:            0,
		Voltage:             230, // Default 230V
		ConnectorID:         1,
		MeterValuesInterval: 30,
		InitialSOC:          20,    // Default 20%
		BatteryCapacity:     60000, // Default 60 kWh
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.OCPPVersion != "1.6" && c.OCPPVersion != "2.0.1" {
		return fmt.Errorf("ocpp_version must be '1.6' or '2.0.1', got '%s'", c.OCPPVersion)
	}

	if c.ChargerID == "" {
		return fmt.Errorf("charger_id is required")
	}

	if c.ServerURL == "" {
		return fmt.Errorf("server_url is required")
	}

	if c.MaxCurrent <= 0 {
		return fmt.Errorf("max_current must be positive")
	}

	if c.MaxPower <= 0 {
		return fmt.Errorf("max_power must be positive")
	}

	if c.MinCurrent < 0 {
		return fmt.Errorf("min_current cannot be negative")
	}

	if c.MinPower < 0 {
		return fmt.Errorf("min_power cannot be negative")
	}

	if c.MinCurrent > c.MaxCurrent {
		return fmt.Errorf("min_current cannot exceed max_current")
	}

	if c.MinPower > c.MaxPower {
		return fmt.Errorf("min_power cannot exceed max_power")
	}

	if c.Voltage <= 0 {
		return fmt.Errorf("voltage must be positive")
	}

	if c.InitialSOC < 0 || c.InitialSOC > 100 {
		return fmt.Errorf("initial_soc must be between 0 and 100")
	}

	if c.BatteryCapacity <= 0 {
		return fmt.Errorf("battery_capacity must be positive")
	}

	return nil
}

// GetTLSConfig returns the tls.Config if TLS is configured
func (c *Config) GetTLSConfig() (*tls.Config, error) {
	if c.TLS == nil {
		return nil, nil
	}

	tlsConfig := &tls.Config{}

	// Skip server certificate verification if requested
	if c.TLS.SkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	// Build certificate pool for trusted certificates
	certPool := x509.NewCertPool()
	hasCerts := false

	// Load CA certificate if provided
	if c.TLS.CAFile != "" {
		caCert, err := os.ReadFile(c.TLS.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}
		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		hasCerts = true
	}

	// Load trusted server certificate if provided (for self-signed certs)
	if c.TLS.ServerCertFile != "" {
		serverCert, err := os.ReadFile(c.TLS.ServerCertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read server certificate: %w", err)
		}
		if !certPool.AppendCertsFromPEM(serverCert) {
			return nil, fmt.Errorf("failed to parse server certificate")
		}
		hasCerts = true
	}

	if hasCerts {
		tlsConfig.RootCAs = certPool
	}

	// Load client certificate and key if provided
	if c.TLS.CertFile != "" && c.TLS.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(c.TLS.CertFile, c.TLS.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

// IsOCPP16 returns true if the configured version is 1.6
func (c *Config) IsOCPP16() bool {
	return c.OCPPVersion == "1.6"
}

// IsOCPP201 returns true if the configured version is 2.0.1
func (c *Config) IsOCPP201() bool {
	return c.OCPPVersion == "2.0.1"
}
