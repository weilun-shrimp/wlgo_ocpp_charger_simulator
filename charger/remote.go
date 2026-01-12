package charger

import (
	"encoding/json"
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

	resp := v16.RemoteStartTransactionResponse{
		Status: "Accepted",
	}

	if err := c.sendCallResult(uniqueId, resp); err != nil {
		log.Printf("Failed to send RemoteStartTransaction response: %v", err)
		return
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err := c.StartTransaction(req.IdTag); err != nil {
			log.Printf("Failed to start transaction: %v", err)
		}
	}()
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

	log.Printf("Received RequestStartTransaction: idToken=%s, evseId=%d", req.IdToken.IdToken, req.EvseId)

	c.mu.Lock()
	c.transactionIdStr = uuid.New().String()
	transactionId := c.transactionIdStr
	c.mu.Unlock()

	resp := v201.RequestStartTransactionResponse{
		Status:        "Accepted",
		TransactionId: transactionId,
	}

	if err := c.sendCallResult(uniqueId, resp); err != nil {
		log.Printf("Failed to send RequestStartTransaction response: %v", err)
		return
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err := c.StartTransaction(req.IdToken.IdToken); err != nil {
			log.Printf("Failed to start transaction: %v", err)
		}
	}()
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

	// Extract current from charging profile (OCPP 1.6 only supports Amperes)
	if req.ChargingProfile != nil && req.ChargingProfile.ChargingSchedule != nil {
		schedule := req.ChargingProfile.ChargingSchedule
		if len(schedule.ChargingSchedulePeriod) > 0 {
			limit := schedule.ChargingSchedulePeriod[0].Limit
			unit := schedule.ChargingRateUnit

			if unit != "A" {
				log.Printf("OCPP 1.6 only supports chargingRateUnit 'A', got: %s", unit)
				status = "Rejected"
			} else {
				if err := c.SetCurrent(limit); err != nil {
					log.Printf("Failed to set current: %v", err)
					status = "Rejected"
				}
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
