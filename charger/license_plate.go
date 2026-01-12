package charger

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v16"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v201"
)

// LicensePlateData represents the license plate data to send
type LicensePlateData struct {
	LicensePlate string `json:"licensePlate"`
	ConnectorId  int    `json:"connectorId"`
}

// SetLicensePlateAndSend sets the license plate locally and sends to server if connected
func (c *Charger) SetLicensePlateAndSend(licensePlate string) error {
	c.mu.Lock()
	c.licensePlate = licensePlate
	isConnected := c.isConnected
	c.mu.Unlock()

	log.Printf("License plate set locally: %s", licensePlate)

	// Send to server if connected
	if isConnected {
		return c.SendLicensePlate(licensePlate)
	}
	return nil
}

// SendLicensePlate sends the license plate to the server via DataTransfer
func (c *Charger) SendLicensePlate(licensePlate string) error {
	if c.config.IsOCPP16() {
		return c.sendLicensePlateV16(licensePlate)
	}
	return c.sendLicensePlateV201(licensePlate)
}

func (c *Charger) sendLicensePlateV16(licensePlate string) error {
	data := LicensePlateData{
		LicensePlate: licensePlate,
		ConnectorId:  c.config.ConnectorID,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal license plate data: %w", err)
	}

	req := v16.DataTransferRequest{
		VendorId:  "LicensePlate",
		MessageId: "EVLicensePlate",
		Data:      string(dataJSON),
	}

	resp, err := c.sendCall(v16.ActionDataTransfer, req)
	if err != nil {
		return fmt.Errorf("DataTransfer failed: %w", err)
	}

	var raw []json.RawMessage
	if err := json.Unmarshal(resp, &raw); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(raw) >= 3 {
		var dtResp v16.DataTransferResponse
		if err := json.Unmarshal(raw[2], &dtResp); err != nil {
			return fmt.Errorf("failed to parse DataTransfer response: %w", err)
		}
		log.Printf("DataTransfer (LicensePlate) response: status=%s", dtResp.Status)
	}

	return nil
}

func (c *Charger) sendLicensePlateV201(licensePlate string) error {
	data := LicensePlateData{
		LicensePlate: licensePlate,
		ConnectorId:  c.config.ConnectorID,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal license plate data: %w", err)
	}

	req := v201.DataTransferRequest{
		VendorId:  "LicensePlate",
		MessageId: "EVLicensePlate",
		Data:      string(dataJSON),
	}

	resp, err := c.sendCall(v201.ActionDataTransfer, req)
	if err != nil {
		return fmt.Errorf("DataTransfer failed: %w", err)
	}

	var raw []json.RawMessage
	if err := json.Unmarshal(resp, &raw); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(raw) >= 3 {
		var dtResp v201.DataTransferResponse
		if err := json.Unmarshal(raw[2], &dtResp); err != nil {
			return fmt.Errorf("failed to parse DataTransfer response: %w", err)
		}
		log.Printf("DataTransfer (LicensePlate) response: status=%s", dtResp.Status)
	}

	return nil
}
