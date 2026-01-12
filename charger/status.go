package charger

import (
	"fmt"
	"log"
	"time"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v16"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v201"
)

// Valid OCPP 1.6 statuses
var validStatusV16 = map[string]bool{
	"Available":     true,
	"Preparing":     true,
	"Charging":      true,
	"SuspendedEVSE": true,
	"SuspendedEV":   true,
	"Finishing":     true,
	"Reserved":      true,
	"Unavailable":   true,
	"Faulted":       true,
}

// Valid OCPP 2.0.1 statuses
var validStatusV201 = map[string]bool{
	"Available":   true,
	"Occupied":    true,
	"Reserved":    true,
	"Unavailable": true,
	"Faulted":     true,
}

// SetStatus updates the local status and sends StatusNotification if connected
func (c *Charger) SetStatus(status string) error {
	// Validate status
	if c.config.IsOCPP16() {
		if !validStatusV16[status] {
			return fmt.Errorf("invalid status for OCPP 1.6: %s", status)
		}
	} else {
		if !validStatusV201[status] {
			return fmt.Errorf("invalid status for OCPP 2.0.1: %s", status)
		}
	}

	c.mu.Lock()
	oldStatus := c.status
	c.status = status
	isConnected := c.isConnected

	// Start meter loop when entering Charging
	shouldStartMeter := status == "Charging" && oldStatus != "Charging"
	// Stop meter loop when leaving Charging
	shouldStopMeter := status != "Charging" && oldStatus == "Charging"

	if shouldStartMeter && c.meterStopCh == nil {
		c.meterStopCh = make(chan struct{})
		c.mu.Unlock()
		go c.StartMeterValuesLoop()
	} else if shouldStopMeter && c.meterStopCh != nil {
		close(c.meterStopCh)
		c.meterStopCh = nil
		c.mu.Unlock()
	} else {
		c.mu.Unlock()
	}

	log.Printf("Status changed to: %s", status)

	// Send to server if connected
	if isConnected {
		return c.StatusNotification(status)
	}
	return nil
}

// StatusNotification sends a StatusNotification request to the server
func (c *Charger) StatusNotification(status string) error {
	c.mu.Lock()
	c.status = status
	c.mu.Unlock()

	if c.config.IsOCPP16() {
		return c.statusNotificationV16(status)
	}
	return c.statusNotificationV201(status)
}

func (c *Charger) statusNotificationV16(status string) error {
	req := v16.StatusNotificationRequest{
		ConnectorId: c.config.ConnectorID,
		ErrorCode:   "NoError",
		Status:      v16.ChargePointStatus(status),
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	_, err := c.sendCall(v16.ActionStatusNotification, req)
	if err != nil {
		return fmt.Errorf("StatusNotification failed: %w", err)
	}

	log.Printf("StatusNotification sent: status=%s", status)
	return nil
}

func (c *Charger) statusNotificationV201(status string) error {
	req := v201.StatusNotificationRequest{
		Timestamp:       time.Now().UTC().Format(time.RFC3339),
		ConnectorStatus: v201.ConnectorStatus(status),
		EvseId:          c.config.ConnectorID,
		ConnectorId:     1,
	}

	_, err := c.sendCall(v201.ActionStatusNotification, req)
	if err != nil {
		return fmt.Errorf("StatusNotification failed: %w", err)
	}

	log.Printf("StatusNotification sent: status=%s", status)
	return nil
}
