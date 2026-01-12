package charger

import (
	"crypto/tls"
	"fmt"
	"log"
	"sync"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/config"
	"github.com/weilun-shrimp/wlgows/client"
	"github.com/weilun-shrimp/wlgows/connection"
)

// Charger represents an OCPP charger simulator
type Charger struct {
	config           *config.Config
	conn             *connection.ClientConn
	tlsConfig        *tls.Config
	mu               sync.RWMutex
	status           string
	transactionId    int
	transactionIdStr string // For OCPP 2.0.1
	meterValue       int
	soc              float64 // State of Charge (0-100%)
	licensePlate     string  // License plate from EV
	idTag            string
	seqNo            int
	isCharging       bool
	isConnected      bool
	stopCh           chan struct{} // Stop channel for connect to server
	meterStopCh      chan struct{} // Stop channel for meter loop
	pendingCalls     map[string]chan []byte
	pendingMu        sync.Mutex
}

// New creates a new Charger instance
func New(cfg *config.Config) (*Charger, error) {
	tlsConfig, err := cfg.GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}

	return &Charger{
		config:       cfg,
		tlsConfig:    tlsConfig,
		status:       cfg.InitialStatus,
		meterValue:   0,
		soc:          cfg.InitialSOC,
		stopCh:       make(chan struct{}),
		pendingCalls: make(map[string]chan []byte),
	}, nil
}

// Connect establishes a WebSocket connection to the server
func (c *Charger) Connect() error {
	c.mu.Lock()
	if c.isConnected {
		c.mu.Unlock()
		return fmt.Errorf("already connected")
	}
	// Create new stop channel for this connection
	c.stopCh = make(chan struct{})
	c.mu.Unlock()

	log.Printf("Connecting to %s...", c.config.ServerURL)

	conn, err := client.Dial(c.config.ServerURL, c.tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	if err := conn.HandShake(); err != nil {
		conn.Close()
		return fmt.Errorf("handshake failed: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.isConnected = true
	c.mu.Unlock()

	log.Printf("Connected successfully")

	go c.receiveMessages()

	return nil
}

// Disconnect disconnects from the server
func (c *Charger) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return
	}

	close(c.stopCh)
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.isConnected = false
	c.isCharging = false
}

// Close closes the connection (for defer)
func (c *Charger) Close() {
	c.Disconnect()
}

// IsConnected returns whether the charger is connected to the server
func (c *Charger) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// IsCharging returns whether the charger is currently charging
func (c *Charger) IsCharging() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isCharging
}

// GetStatus returns the current status
func (c *Charger) GetStatus() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

// GetSOC returns the current State of Charge
func (c *Charger) GetSOC() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.soc
}

// SetSOC sets the State of Charge (0-100)
func (c *Charger) SetSOC(soc float64) error {
	if soc < 0 || soc > 100 {
		return fmt.Errorf("SOC must be between 0 and 100")
	}
	c.mu.Lock()
	c.soc = soc
	c.mu.Unlock()
	return nil
}

// GetLicensePlate returns the current license plate
func (c *Charger) GetLicensePlate() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.licensePlate
}

// SetLicensePlate sets the license plate (simulates EV sending plate when plugged in)
func (c *Charger) SetLicensePlate(plate string) {
	c.mu.Lock()
	c.licensePlate = plate
	c.mu.Unlock()
}

// Plugin simulates car plugging in
func (c *Charger) Plugin() error {
	c.mu.RLock()
	status := c.status
	c.mu.RUnlock()

	if status != "Available" {
		return fmt.Errorf("cannot plug in: status must be Available (current: %s)", status)
	}

	return c.SetStatus("Preparing")
}

// Unplug simulates car unplugging - stops all background tasks and resets state
func (c *Charger) Unplug() error {
	c.mu.Lock()
	// Stop meter loop if running
	if c.meterStopCh != nil {
		close(c.meterStopCh)
		c.meterStopCh = nil
	}
	// Reset charging state
	c.isCharging = false
	c.licensePlate = ""
	c.idTag = ""
	c.soc = c.config.InitialSOC
	c.meterValue = 0
	c.mu.Unlock()

	return c.SetStatus("Available")
}
