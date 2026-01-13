package charger

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v16"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v201"
)

// BootNotification sends a BootNotification request
func (c *Charger) BootNotification() error {
	if c.config.IsOCPP16() {
		return c.bootNotificationV16()
	}
	return c.bootNotificationV201()
}

func (c *Charger) bootNotificationV16() error {
	req := v16.BootNotificationRequest{
		ChargePointVendor:       "Simulator",
		ChargePointModel:        "WLGO-SIM-1",
		ChargePointSerialNumber: c.config.ChargerID,
		FirmwareVersion:         "1.0.0",
	}

	resp, err := c.sendCall(v16.ActionBootNotification, req)
	if err != nil {
		return fmt.Errorf("BootNotification failed: %w", err)
	}

	var raw []json.RawMessage
	if err := json.Unmarshal(resp, &raw); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(raw) >= 3 {
		var bootResp v16.BootNotificationResponse
		if err := json.Unmarshal(raw[2], &bootResp); err != nil {
			return fmt.Errorf("failed to parse BootNotification response: %w", err)
		}
		log.Printf("BootNotification response: status=%s, interval=%d", bootResp.Status, bootResp.Interval)

		// Use server-provided interval if available
		if bootResp.Interval > 0 {
			c.SetHeartbeatInterval(bootResp.Interval)
		}
		// Start heartbeat loop
		go c.StartHeartbeatLoop()
	}

	return nil
}

func (c *Charger) bootNotificationV201() error {
	req := v201.BootNotificationRequest{
		Reason: "PowerUp",
		ChargingStation: v201.ChargingStation{
			VendorName:      "Simulator",
			Model:           "WLGO-SIM-2",
			SerialNumber:    c.config.ChargerID,
			FirmwareVersion: "2.0.0",
		},
	}

	resp, err := c.sendCall(v201.ActionBootNotification, req)
	if err != nil {
		return fmt.Errorf("BootNotification failed: %w", err)
	}

	var raw []json.RawMessage
	if err := json.Unmarshal(resp, &raw); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(raw) >= 3 {
		var bootResp v201.BootNotificationResponse
		if err := json.Unmarshal(raw[2], &bootResp); err != nil {
			return fmt.Errorf("failed to parse BootNotification response: %w", err)
		}
		log.Printf("BootNotification response: status=%s, interval=%d", bootResp.Status, bootResp.Interval)

		// Use server-provided interval if available
		if bootResp.Interval > 0 {
			c.SetHeartbeatInterval(bootResp.Interval)
		}
		// Start heartbeat loop
		go c.StartHeartbeatLoop()
	}

	return nil
}
