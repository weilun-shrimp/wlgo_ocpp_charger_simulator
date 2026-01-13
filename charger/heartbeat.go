package charger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v16"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v201"
)

// Heartbeat sends a Heartbeat request to the server
func (c *Charger) Heartbeat() error {
	c.mu.RLock()
	isConnected := c.isConnected
	c.mu.RUnlock()

	if !isConnected {
		return fmt.Errorf("not connected to server")
	}

	if c.config.IsOCPP16() {
		return c.heartbeatV16()
	}
	return c.heartbeatV201()
}

func (c *Charger) heartbeatV16() error {
	req := v16.HeartbeatRequest{}

	resp, err := c.sendCall(v16.ActionHeartbeat, req)
	if err != nil {
		return fmt.Errorf("Heartbeat failed: %w", err)
	}

	var raw []json.RawMessage
	if err := json.Unmarshal(resp, &raw); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(raw) >= 3 {
		var heartbeatResp v16.HeartbeatResponse
		if err := json.Unmarshal(raw[2], &heartbeatResp); err != nil {
			return fmt.Errorf("failed to parse Heartbeat response: %w", err)
		}
		log.Printf("Heartbeat response: currentTime=%s", heartbeatResp.CurrentTime)
	}

	return nil
}

func (c *Charger) heartbeatV201() error {
	req := v201.HeartbeatRequest{}

	resp, err := c.sendCall(v201.ActionHeartbeat, req)
	if err != nil {
		return fmt.Errorf("Heartbeat failed: %w", err)
	}

	var raw []json.RawMessage
	if err := json.Unmarshal(resp, &raw); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(raw) >= 3 {
		var heartbeatResp v201.HeartbeatResponse
		if err := json.Unmarshal(raw[2], &heartbeatResp); err != nil {
			return fmt.Errorf("failed to parse Heartbeat response: %w", err)
		}
		log.Printf("Heartbeat response: currentTime=%s", heartbeatResp.CurrentTime)
	}

	return nil
}

// StartHeartbeatLoop starts the heartbeat loop with the configured interval
func (c *Charger) StartHeartbeatLoop() {
	c.mu.Lock()
	interval := c.heartbeatInterval
	if interval <= 0 {
		c.mu.Unlock()
		log.Printf("Heartbeat disabled (interval=%d)", interval)
		return
	}
	c.heartbeatStopCh = make(chan struct{})
	stopCh := c.heartbeatStopCh
	c.mu.Unlock()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	log.Printf("Heartbeat loop started (interval=%ds)", interval)

	for {
		select {
		case <-stopCh:
			log.Printf("Heartbeat loop stopped")
			return
		case <-ticker.C:
			if err := c.Heartbeat(); err != nil {
				log.Printf("Heartbeat error: %v", err)
			}
		}
	}
}

// StopHeartbeatLoop stops the heartbeat loop
func (c *Charger) StopHeartbeatLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.heartbeatStopCh != nil {
		close(c.heartbeatStopCh)
		c.heartbeatStopCh = nil
	}
}

// SetHeartbeatInterval updates the heartbeat interval (e.g., from BootNotification response)
func (c *Charger) SetHeartbeatInterval(interval int) {
	c.mu.Lock()
	c.heartbeatInterval = interval
	c.mu.Unlock()
	log.Printf("Heartbeat interval set to %d seconds", interval)
}
