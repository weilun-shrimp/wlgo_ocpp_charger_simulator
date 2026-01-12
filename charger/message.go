package charger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v16"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v201"
)

// receiveMessages handles incoming messages from the server
func (c *Charger) receiveMessages() {
	defer c.Disconnect()

	for {
		select {
		case <-c.stopCh:
			return
		default:
			msg, err := c.conn.GetNextMsg()
			if err != nil {
				if err == io.EOF {
					log.Printf("Server closed connection (EOF)")
				} else {
					log.Printf("Error receiving message: %v", err)
				}
				return
			}

			data := msg.GetStr()
			log.Printf("Received: %s", data)

			go c.handleMessage([]byte(data))
		}
	}
}

// handleMessage processes incoming OCPP messages
func (c *Charger) handleMessage(data []byte) {
	if c.config.IsOCPP16() {
		c.handleMessageV16(data)
	} else {
		c.handleMessageV201(data)
	}
}

// handleMessageV16 processes OCPP 1.6 messages
func (c *Charger) handleMessageV16(data []byte) {
	messageType, uniqueId, payload, action, err := v16.ParseMessage(data)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	switch messageType {
	case v16.MessageTypeCall:
		c.handleCallV16(uniqueId, action, payload)
	case v16.MessageTypeCallResult:
		c.handleCallResult(uniqueId, data)
	case v16.MessageTypeCallError:
		log.Printf("Received CallError for %s: %s", uniqueId, string(payload))
		c.handleCallResult(uniqueId, data)
	}
}

// handleMessageV201 processes OCPP 2.0.1 messages
func (c *Charger) handleMessageV201(data []byte) {
	messageType, uniqueId, payload, action, err := v201.ParseMessage(data)
	if err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	switch messageType {
	case v201.MessageTypeCall:
		c.handleCallV201(uniqueId, action, payload)
	case v201.MessageTypeCallResult:
		c.handleCallResult(uniqueId, data)
	case v201.MessageTypeCallError:
		log.Printf("Received CallError for %s: %s", uniqueId, string(payload))
		c.handleCallResult(uniqueId, data)
	}
}

// handleCallResult notifies waiting goroutines about the response
func (c *Charger) handleCallResult(uniqueId string, data []byte) {
	c.pendingMu.Lock()
	ch, ok := c.pendingCalls[uniqueId]
	if ok {
		delete(c.pendingCalls, uniqueId)
	}
	c.pendingMu.Unlock()

	if ok {
		ch <- data
	}
}

// handleCallV16 processes incoming OCPP 1.6 Call messages
func (c *Charger) handleCallV16(uniqueId, action string, payload json.RawMessage) {
	switch action {
	case v16.ActionRemoteStartTransaction:
		c.handleRemoteStartTransactionV16(uniqueId, payload)
	case v16.ActionRemoteStopTransaction:
		c.handleRemoteStopTransactionV16(uniqueId, payload)
	case v16.ActionSetChargingProfile:
		c.handleSetChargingProfileV16(uniqueId, payload)
	default:
		log.Printf("Unknown action: %s", action)
	}
}

// handleCallV201 processes incoming OCPP 2.0.1 Call messages
func (c *Charger) handleCallV201(uniqueId, action string, payload json.RawMessage) {
	switch action {
	case v201.ActionRequestStartTransaction:
		c.handleRequestStartTransactionV201(uniqueId, payload)
	case v201.ActionRequestStopTransaction:
		c.handleRequestStopTransactionV201(uniqueId, payload)
	case v201.ActionSetChargingProfile:
		c.handleSetChargingProfileV201(uniqueId, payload)
	default:
		log.Printf("Unknown action: %s", action)
	}
}

// sendCall sends a Call message and waits for response
func (c *Charger) sendCall(action string, payload interface{}) ([]byte, error) {
	uniqueId := uuid.New().String()

	var data []byte
	var err error

	if c.config.IsOCPP16() {
		data, err = v16.MarshalCall(uniqueId, action, payload)
	} else {
		data, err = v201.MarshalCall(uniqueId, action, payload)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	respCh := make(chan []byte, 1)
	c.pendingMu.Lock()
	c.pendingCalls[uniqueId] = respCh
	c.pendingMu.Unlock()

	log.Printf("Sending: %s", string(data))
	c.conn.SendText(data)

	select {
	case resp := <-respCh:
		return resp, nil
	case <-time.After(30 * time.Second):
		c.pendingMu.Lock()
		delete(c.pendingCalls, uniqueId)
		c.pendingMu.Unlock()
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

// sendCallResult sends a CallResult message
func (c *Charger) sendCallResult(uniqueId string, payload interface{}) error {
	var data []byte
	var err error

	if c.config.IsOCPP16() {
		data, err = v16.MarshalCallResult(uniqueId, payload)
	} else {
		data, err = v201.MarshalCallResult(uniqueId, payload)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	log.Printf("Sending: %s", string(data))
	c.conn.SendText(data)
	return nil
}
