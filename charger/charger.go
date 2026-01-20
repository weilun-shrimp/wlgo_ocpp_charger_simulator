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
	current           float64       // Current limit in Amperes (between MinCurrent and MaxCurrent)
	power             float64       // Power limit in Watts (between MinPower and MaxPower)
	stopCh            chan struct{} // Stop channel for connect to server
	meterStopCh       chan struct{} // Stop channel for meter loop
	heartbeatInterval int           // Heartbeat interval in seconds (from config or server)
	heartbeatStopCh   chan struct{} // Stop channel for heartbeat loop
	pendingCalls      map[string]chan []byte
	pendingMu         sync.Mutex
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
		current:      cfg.MaxCurrent, // Default to max current
		power:        cfg.MaxPower,   // Default to max power
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

	// Stop heartbeat loop
	if c.heartbeatStopCh != nil {
		close(c.heartbeatStopCh)
		c.heartbeatStopCh = nil
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

// GetCurrent returns the current in Amperes
func (c *Charger) GetCurrent() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.current
}

// SetCurrent sets the current in Amperes (bounded by MinCurrent and MaxCurrent)
// Setting current to 0 will suspend charging (SuspendedEVSE)
// Setting current > 0 from SuspendedEVSE will resume charging
func (c *Charger) SetCurrent(current float64) error {
	// Allow 0 for suspend, otherwise check MinCurrent
	if current != 0 && current < c.config.MinCurrent {
		return fmt.Errorf("current %.1fA is below minimum %.1fA", current, c.config.MinCurrent)
	}
	if current > c.config.MaxCurrent {
		return fmt.Errorf("current %.1fA exceeds maximum %.1fA", current, c.config.MaxCurrent)
	}

	c.mu.Lock()
	oldCurrent := c.current
	c.current = current
	status := c.status
	c.mu.Unlock()

	log.Printf("Current set to %.1f A", current)

	// Handle status transitions per OCPP spec:
	// - Charging -> SuspendedEVSE when EVSE sets current to 0
	// - SuspendedEVSE -> Charging when EVSE restores current > 0
	if current == 0 && oldCurrent > 0 && status == "Charging" {
		return c.SetStatus("SuspendedEVSE")
	} else if current > 0 && oldCurrent == 0 && status == "SuspendedEVSE" {
		return c.SetStatus("Charging")
	}

	return nil
}

// GetPower returns the power in Watts
func (c *Charger) GetPower() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.power
}

// SetPower sets the power in Watts (bounded by MinPower and MaxPower)
// Setting power to 0 will suspend charging (SuspendedEVSE)
// Setting power > 0 from SuspendedEVSE will resume charging
func (c *Charger) SetPower(power float64) error {
	// Allow 0 for suspend, otherwise check MinPower
	if power != 0 && power < c.config.MinPower {
		return fmt.Errorf("power %.1fW is below minimum %.1fW", power, c.config.MinPower)
	}
	if power > c.config.MaxPower {
		return fmt.Errorf("power %.1fW exceeds maximum %.1fW", power, c.config.MaxPower)
	}

	c.mu.Lock()
	oldPower := c.power
	c.power = power
	// Also update current based on power (I = P / V)
	if power > 0 {
		c.current = power / c.config.Voltage
	} else {
		c.current = 0
	}
	status := c.status
	c.mu.Unlock()

	log.Printf("Power set to %.1f W (current: %.1f A)", power, power/c.config.Voltage)

	// Handle status transitions per OCPP spec:
	// - Charging -> SuspendedEVSE when EVSE sets power to 0
	// - SuspendedEVSE -> Charging when EVSE restores power > 0
	if power == 0 && oldPower > 0 && status == "Charging" {
		return c.SetStatus("SuspendedEVSE")
	} else if power > 0 && oldPower == 0 && status == "SuspendedEVSE" {
		return c.SetStatus("Charging")
	}

	return nil
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
