package charger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v16"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v201"
)

// handleRemoteStartTransactionV16 handles RemoteStartTransaction from server
func (c *Charger) handleRemoteStartTransactionV16(uniqueId string, payload json.RawMessage) {
	var req v16.RemoteStartTransactionRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("Failed to parse RemoteStartTransaction: %v", err)
		return
	}

	log.Printf("Received RemoteStartTransaction: idTag=%s, connectorId=%d", req.IdTag, req.ConnectorId)

	c.mu.Lock()
	status := c.status
	var respStatus string

	switch status {
	case "Available":
		// Accept and store pending authorization - will auto-start when cable plugged in
		respStatus = "Accepted"
		c.pendingRemoteStartIdTag = req.IdTag
		log.Printf("RemoteStartTransaction accepted: waiting for cable to be plugged in")
	case "Preparing":
		// Cable already plugged in - accept and start immediately
		respStatus = "Accepted"
		c.pendingRemoteStartIdTag = "" // Clear any pending
	default:
		// Reject if charging, finishing, or other states
		respStatus = "Rejected"
		log.Printf("RemoteStartTransaction rejected: invalid status (current: %s)", status)
	}
	c.mu.Unlock()

	resp := v16.RemoteStartTransactionResponse{
		Status: respStatus,
	}

	if err := c.sendCallResult(uniqueId, resp); err != nil {
		log.Printf("Failed to send RemoteStartTransaction response: %v", err)
		return
	}

	// If cable is already plugged in (Preparing), start transaction immediately
	if respStatus == "Accepted" && status == "Preparing" {
		go func() {
			time.Sleep(1 * time.Second)
			if err := c.StartTransaction(req.IdTag); err != nil {
				log.Printf("Failed to start transaction: %v", err)
			}
		}()
	}
}

// handleRemoteStopTransactionV16 handles RemoteStopTransaction from server
func (c *Charger) handleRemoteStopTransactionV16(uniqueId string, payload json.RawMessage) {
	var req v16.RemoteStopTransactionRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("Failed to parse RemoteStopTransaction: %v", err)
		return
	}

	log.Printf("Received RemoteStopTransaction: transactionId=%d", req.TransactionId)

	c.mu.RLock()
	currentTransactionId := c.transactionId
	isCharging := c.isCharging
	c.mu.RUnlock()

	var status string
	if isCharging && currentTransactionId == req.TransactionId {
		status = "Accepted"
	} else {
		status = "Rejected"
	}

	resp := v16.RemoteStopTransactionResponse{
		Status: status,
	}

	if err := c.sendCallResult(uniqueId, resp); err != nil {
		log.Printf("Failed to send RemoteStopTransaction response: %v", err)
		return
	}

	if status == "Accepted" {
		go func() {
			time.Sleep(1 * time.Second)
			if err := c.StopTransaction("Remote"); err != nil {
				log.Printf("Failed to stop transaction: %v", err)
			}
		}()
	}
}

// handleRequestStartTransactionV201 handles RequestStartTransaction from server
func (c *Charger) handleRequestStartTransactionV201(uniqueId string, payload json.RawMessage) {
	var req v201.RequestStartTransactionRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("Failed to parse RequestStartTransaction: %v", err)
		return
	}

	log.Printf("Received RequestStartTransaction: idToken=%s, evseId=%d, remoteStartId=%d", req.IdToken.IdToken, req.EvseId, req.RemoteStartId)

	c.mu.Lock()
	status := c.status
	var respStatus string
	var transactionId string
	var statusInfo *v201.StatusInfo

	switch status {
	case "Available":
		// Accept and store pending authorization - will auto-start when cable plugged in
		respStatus = "Accepted"
		c.pendingRemoteStartIdTag = req.IdToken.IdToken
		c.pendingRemoteStartId = req.RemoteStartId
		// Generate transaction ID now for the response
		c.transactionIdStr = uuid.New().String()
		transactionId = c.transactionIdStr
		log.Printf("RequestStartTransaction accepted: waiting for cable to be plugged in")
	case "Occupied":
		// Cable already plugged in - accept and start immediately
		respStatus = "Accepted"
		c.pendingRemoteStartIdTag = "" // Clear any pending
		c.pendingRemoteStartId = 0
		c.transactionIdStr = uuid.New().String()
		transactionId = c.transactionIdStr
	default:
		// Reject if charging or other states
		respStatus = "Rejected"
		statusInfo = &v201.StatusInfo{
			ReasonCode:     "Occupied",
			AdditionalInfo: fmt.Sprintf("Charger is busy, current status: %s", status),
		}
		log.Printf("RequestStartTransaction rejected: invalid status (current: %s)", status)
	}
	c.mu.Unlock()

	resp := v201.RequestStartTransactionResponse{
		Status:        respStatus,
		TransactionId: transactionId,
		StatusInfo:    statusInfo,
	}

	if err := c.sendCallResult(uniqueId, resp); err != nil {
		log.Printf("Failed to send RequestStartTransaction response: %v", err)
		return
	}

	// If cable is already plugged in (Occupied), start transaction immediately
	if respStatus == "Accepted" && status == "Occupied" {
		go func() {
			time.Sleep(1 * time.Second)
			if err := c.StartTransaction(req.IdToken.IdToken); err != nil {
				log.Printf("Failed to start transaction: %v", err)
			}
		}()
	}
}

// handleRequestStopTransactionV201 handles RequestStopTransaction from server
func (c *Charger) handleRequestStopTransactionV201(uniqueId string, payload json.RawMessage) {
	var req v201.RequestStopTransactionRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("Failed to parse RequestStopTransaction: %v", err)
		return
	}

	log.Printf("Received RequestStopTransaction: transactionId=%s", req.TransactionId)

	c.mu.RLock()
	currentTransactionId := c.transactionIdStr
	isCharging := c.isCharging
	c.mu.RUnlock()

	var status string
	if isCharging && currentTransactionId == req.TransactionId {
		status = "Accepted"
	} else {
		status = "Rejected"
	}

	resp := v201.RequestStopTransactionResponse{
		Status: status,
	}

	if err := c.sendCallResult(uniqueId, resp); err != nil {
		log.Printf("Failed to send RequestStopTransaction response: %v", err)
		return
	}

	if status == "Accepted" {
		go func() {
			time.Sleep(1 * time.Second)
			if err := c.StopTransaction("Remote"); err != nil {
				log.Printf("Failed to stop transaction: %v", err)
			}
		}()
	}
}

// handleSetChargingProfileV16 handles SetChargingProfile from server (OCPP 1.6)
func (c *Charger) handleSetChargingProfileV16(uniqueId string, payload json.RawMessage) {
	var req v16.SetChargingProfileRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("Failed to parse SetChargingProfile: %v", err)
		return
	}

	log.Printf("Received SetChargingProfile: connectorId=%d", req.ConnectorId)

	status := "Accepted"

	// Extract limit from charging profile
	if req.ChargingProfile != nil && req.ChargingProfile.ChargingSchedule != nil {
		schedule := req.ChargingProfile.ChargingSchedule
		if len(schedule.ChargingSchedulePeriod) > 0 {
			limit := schedule.ChargingSchedulePeriod[0].Limit
			unit := schedule.ChargingRateUnit

			if unit == "A" {
				if err := c.SetCurrent(limit); err != nil {
					log.Printf("Failed to set current: %v", err)
					status = "Rejected"
				}
			} else if unit == "W" {
				if err := c.SetPower(limit); err != nil {
					log.Printf("Failed to set power: %v", err)
					status = "Rejected"
				}
			} else {
				log.Printf("Unknown chargingRateUnit: %s", unit)
				status = "Rejected"
			}
		}
	}

	resp := v16.SetChargingProfileResponse{
		Status: status,
	}

	if err := c.sendCallResult(uniqueId, resp); err != nil {
		log.Printf("Failed to send SetChargingProfile response: %v", err)
	}
}

// handleSetChargingProfileV201 handles SetChargingProfile from server (OCPP 2.0.1)
func (c *Charger) handleSetChargingProfileV201(uniqueId string, payload json.RawMessage) {
	var req v201.SetChargingProfileRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("Failed to parse SetChargingProfile: %v", err)
		return
	}

	log.Printf("Received SetChargingProfile: evseId=%d", req.EvseId)

	status := "Accepted"

	// Extract current limit from charging profile
	if req.ChargingProfile != nil && len(req.ChargingProfile.ChargingSchedule) > 0 {
		schedule := req.ChargingProfile.ChargingSchedule[0]
		if len(schedule.ChargingSchedulePeriod) > 0 {
			limit := schedule.ChargingSchedulePeriod[0].Limit
			unit := schedule.ChargingRateUnit

			var currentAmps float64
			if unit == "A" {
				currentAmps = limit
			} else if unit == "W" {
				// Convert power to current
				currentAmps = limit / c.config.Voltage
			} else {
				log.Printf("Unknown charging rate unit: %s", unit)
				status = "Rejected"
			}

			if status == "Accepted" {
				if err := c.SetCurrent(currentAmps); err != nil {
					log.Printf("Failed to set current: %v", err)
					status = "Rejected"
				}
			}
		}
	}

	resp := v201.SetChargingProfileResponse{
		Status: status,
	}

	if err := c.sendCallResult(uniqueId, resp); err != nil {
		log.Printf("Failed to send SetChargingProfile response: %v", err)
	}
}
