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

// StartTransaction starts a transaction locally and sends to server if connected
func (c *Charger) StartTransaction(idTag string) error {
	c.mu.Lock()
	// OCPP 1.6 requires "Preparing", OCPP 2.0.1 requires "Occupied"
	requiredStatus := "Preparing"
	if !c.config.IsOCPP16() {
		requiredStatus = "Occupied"
	}
	if c.status != requiredStatus {
		c.mu.Unlock()
		return fmt.Errorf("cannot start transaction: status must be %s (current: %s)", requiredStatus, c.status)
	}
	c.idTag = idTag
	c.meterValue = 0
	c.seqNo = 0
	c.isCharging = true
	isConnected := c.isConnected

	// For OCPP 2.0.1, start meter loop here since we don't change status to "Charging"
	shouldStartMeter := !c.config.IsOCPP16() && c.meterStopCh == nil
	if shouldStartMeter {
		c.meterStopCh = make(chan struct{})
	}
	c.mu.Unlock()

	log.Printf("Transaction started locally: idTag=%s", idTag)

	// Start meter loop for OCPP 2.0.1 (OCPP 1.6 starts it via SetStatus("Charging"))
	if shouldStartMeter {
		go c.StartMeterValuesLoop()
	}

	// Update status locally (and send if connected)
	// OCPP 1.6: Status changes to "Charging" (this also starts the meter loop)
	// OCPP 2.0.1: Status stays "Occupied" (charging state is in TransactionEvent)
	if c.config.IsOCPP16() {
		c.SetStatus("Charging")
	}

	// Send to server if connected
	if isConnected {
		if c.config.IsOCPP16() {
			return c.sendStartTransactionV16(idTag)
		}
		return c.sendStartTransactionV201(idTag)
	}
	return nil
}

func (c *Charger) sendStartTransactionV16(idTag string) error {
	req := v16.StartTransactionRequest{
		ConnectorId: c.config.ConnectorID,
		IdTag:       idTag,
		MeterStart:  c.meterValue,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	resp, err := c.sendCall(v16.ActionStartTransaction, req)
	if err != nil {
		return fmt.Errorf("StartTransaction failed: %w", err)
	}

	var raw []json.RawMessage
	if err := json.Unmarshal(resp, &raw); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(raw) >= 3 {
		var startResp v16.StartTransactionResponse
		if err := json.Unmarshal(raw[2], &startResp); err != nil {
			return fmt.Errorf("failed to parse StartTransaction response: %w", err)
		}

		c.mu.Lock()
		c.transactionId = startResp.TransactionId
		c.mu.Unlock()

		log.Printf("StartTransaction response: transactionId=%d, status=%s", startResp.TransactionId, startResp.IdTagInfo.Status)
	}

	return nil
}

func (c *Charger) sendStartTransactionV201(idTag string) error {
	c.mu.Lock()
	c.transactionIdStr = uuid.New().String()
	transactionIdStr := c.transactionIdStr
	c.mu.Unlock()

	req := v201.TransactionEventRequest{
		EventType:     v201.TransactionEventStarted,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		TriggerReason: v201.TriggerReasonAuthorized,
		SeqNo:         0,
		TransactionInfo: v201.Transaction{
			TransactionId: transactionIdStr,
			ChargingState: v201.ChargingStateCharging,
		},
		Evse: &v201.EVSE{
			Id:          c.config.ConnectorID,
			ConnectorId: 1,
		},
		IdToken: &v201.IdToken{
			IdToken: idTag,
			Type:    "ISO14443",
		},
	}

	resp, err := c.sendCall(v201.ActionTransactionEvent, req)
	if err != nil {
		return fmt.Errorf("TransactionEvent (Started) failed: %w", err)
	}

	log.Printf("TransactionEvent (Started) sent: transactionId=%s", transactionIdStr)

	var raw []json.RawMessage
	if err := json.Unmarshal(resp, &raw); err == nil && len(raw) >= 3 {
		var eventResp v201.TransactionEventResponse
		if err := json.Unmarshal(raw[2], &eventResp); err == nil {
			log.Printf("TransactionEvent response received")
		}
	}

	return nil
}

// StopTransaction stops a transaction locally and sends to server if connected
func (c *Charger) StopTransaction(reason string) error {
	c.mu.Lock()
	c.isCharging = false
	meterValue := c.meterValue
	transactionId := c.transactionId
	transactionIdStr := c.transactionIdStr
	idTag := c.idTag
	c.seqNo++
	seqNo := c.seqNo
	isConnected := c.isConnected

	// For OCPP 2.0.1, stop meter loop here since we don't change status from "Charging"
	if !c.config.IsOCPP16() && c.meterStopCh != nil {
		close(c.meterStopCh)
		c.meterStopCh = nil
	}
	c.mu.Unlock()

	log.Printf("Transaction stopped locally: reason=%s", reason)

	// Update status locally (and send if connected)
	// OCPP 1.6: Status changes to "Finishing" (this also stops the meter loop)
	// OCPP 2.0.1: Status stays "Occupied" (cable still connected)
	if c.config.IsOCPP16() {
		c.SetStatus("Finishing")
	}

	// Send to server if connected
	if isConnected {
		if c.config.IsOCPP16() {
			return c.sendStopTransactionV16(meterValue, transactionId, idTag, reason)
		}
		return c.sendStopTransactionV201(meterValue, transactionIdStr, seqNo, reason)
	}
	return nil
}

func (c *Charger) sendStopTransactionV16(meterValue, transactionId int, idTag, reason string) error {
	req := v16.StopTransactionRequest{
		IdTag:         idTag,
		MeterStop:     meterValue,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		TransactionId: transactionId,
		Reason:        reason,
	}

	_, err := c.sendCall(v16.ActionStopTransaction, req)
	if err != nil {
		return fmt.Errorf("StopTransaction failed: %w", err)
	}

	log.Printf("StopTransaction sent: transactionId=%d, meterStop=%d, reason=%s", transactionId, meterValue, reason)

	return nil
}

func (c *Charger) sendStopTransactionV201(meterValue int, transactionIdStr string, seqNo int, reason string) error {
	req := v201.TransactionEventRequest{
		EventType:     v201.TransactionEventEnded,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		TriggerReason: v201.TriggerReasonStopAuthorized,
		SeqNo:         seqNo,
		TransactionInfo: v201.Transaction{
			TransactionId: transactionIdStr,
			ChargingState: v201.ChargingStateIdle,
			StoppedReason: reason,
		},
		MeterValue: []v201.MeterValue{
			{
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				SampledValue: []v201.SampledValue{
					{
						Value:     float64(meterValue),
						Context:   "Transaction.End",
						Measurand: "Energy.Active.Import.Register",
						UnitOfMeasure: &v201.UnitOfMeasure{
							Unit: "Wh",
						},
					},
				},
			},
		},
	}

	_, err := c.sendCall(v201.ActionTransactionEvent, req)
	if err != nil {
		return fmt.Errorf("TransactionEvent (Ended) failed: %w", err)
	}

	log.Printf("TransactionEvent (Ended) sent: transactionId=%s, meterStop=%d, reason=%s", transactionIdStr, meterValue, reason)

	return nil
}
