package v16

import (
	"encoding/json"
	"fmt"
)

// OCPP 1.6 Message Types
const (
	MessageTypeCall       = 2 // Request
	MessageTypeCallResult = 3 // Response
	MessageTypeCallError  = 4 // Error
)

// OCPP 1.6 Actions
const (
	ActionBootNotification       = "BootNotification"
	ActionStatusNotification     = "StatusNotification"
	ActionStartTransaction       = "StartTransaction"
	ActionStopTransaction        = "StopTransaction"
	ActionMeterValues            = "MeterValues"
	ActionRemoteStartTransaction = "RemoteStartTransaction"
	ActionRemoteStopTransaction  = "RemoteStopTransaction"
	ActionHeartbeat              = "Heartbeat"
	ActionDataTransfer           = "DataTransfer"
)

// ChargePointStatus represents the status of a charge point
type ChargePointStatus string

const (
	StatusAvailable     ChargePointStatus = "Available"
	StatusPreparing     ChargePointStatus = "Preparing"
	StatusCharging      ChargePointStatus = "Charging"
	StatusSuspendedEVSE ChargePointStatus = "SuspendedEVSE"
	StatusSuspendedEV   ChargePointStatus = "SuspendedEV"
	StatusFinishing     ChargePointStatus = "Finishing"
	StatusReserved      ChargePointStatus = "Reserved"
	StatusUnavailable   ChargePointStatus = "Unavailable"
	StatusFaulted       ChargePointStatus = "Faulted"
)

// RegistrationStatus represents the registration status in BootNotification response
type RegistrationStatus string

const (
	RegistrationAccepted RegistrationStatus = "Accepted"
	RegistrationPending  RegistrationStatus = "Pending"
	RegistrationRejected RegistrationStatus = "Rejected"
)

// BootNotificationRequest is the request for BootNotification
type BootNotificationRequest struct {
	ChargePointVendor       string `json:"chargePointVendor"`
	ChargePointModel        string `json:"chargePointModel"`
	ChargePointSerialNumber string `json:"chargePointSerialNumber,omitempty"`
	ChargeBoxSerialNumber   string `json:"chargeBoxSerialNumber,omitempty"`
	FirmwareVersion         string `json:"firmwareVersion,omitempty"`
	Iccid                   string `json:"iccid,omitempty"`
	Imsi                    string `json:"imsi,omitempty"`
	MeterType               string `json:"meterType,omitempty"`
	MeterSerialNumber       string `json:"meterSerialNumber,omitempty"`
}

// BootNotificationResponse is the response for BootNotification
type BootNotificationResponse struct {
	Status      RegistrationStatus `json:"status"`
	CurrentTime string             `json:"currentTime"`
	Interval    int                `json:"interval"`
}

// StatusNotificationRequest is the request for StatusNotification
type StatusNotificationRequest struct {
	ConnectorId     int               `json:"connectorId"`
	ErrorCode       string            `json:"errorCode"`
	Status          ChargePointStatus `json:"status"`
	Timestamp       string            `json:"timestamp,omitempty"`
	Info            string            `json:"info,omitempty"`
	VendorId        string            `json:"vendorId,omitempty"`
	VendorErrorCode string            `json:"vendorErrorCode,omitempty"`
}

// StatusNotificationResponse is the response for StatusNotification
type StatusNotificationResponse struct{}

// StartTransactionRequest is the request for StartTransaction
type StartTransactionRequest struct {
	ConnectorId   int    `json:"connectorId"`
	IdTag         string `json:"idTag"`
	MeterStart    int    `json:"meterStart"`
	Timestamp     string `json:"timestamp"`
	ReservationId int    `json:"reservationId,omitempty"`
}

// StartTransactionResponse is the response for StartTransaction
type StartTransactionResponse struct {
	IdTagInfo     IdTagInfo `json:"idTagInfo"`
	TransactionId int       `json:"transactionId"`
}

// IdTagInfo contains authorization information
type IdTagInfo struct {
	Status      string `json:"status"`
	ExpiryDate  string `json:"expiryDate,omitempty"`
	ParentIdTag string `json:"parentIdTag,omitempty"`
}

// StopTransactionRequest is the request for StopTransaction
type StopTransactionRequest struct {
	IdTag           string            `json:"idTag,omitempty"`
	MeterStop       int               `json:"meterStop"`
	Timestamp       string            `json:"timestamp"`
	TransactionId   int               `json:"transactionId"`
	Reason          string            `json:"reason,omitempty"`
	TransactionData []MeterValueEntry `json:"transactionData,omitempty"`
}

// StopTransactionResponse is the response for StopTransaction
type StopTransactionResponse struct {
	IdTagInfo *IdTagInfo `json:"idTagInfo,omitempty"`
}

// MeterValuesRequest is the request for MeterValues
type MeterValuesRequest struct {
	ConnectorId   int               `json:"connectorId"`
	TransactionId int               `json:"transactionId,omitempty"`
	MeterValue    []MeterValueEntry `json:"meterValue"`
}

// MeterValueEntry represents a meter value entry
type MeterValueEntry struct {
	Timestamp    string         `json:"timestamp"`
	SampledValue []SampledValue `json:"sampledValue"`
}

// SampledValue represents a sampled value
type SampledValue struct {
	Value     string `json:"value"`
	Context   string `json:"context,omitempty"`
	Format    string `json:"format,omitempty"`
	Measurand string `json:"measurand,omitempty"`
	Phase     string `json:"phase,omitempty"`
	Location  string `json:"location,omitempty"`
	Unit      string `json:"unit,omitempty"`
}

// MeterValuesResponse is the response for MeterValues
type MeterValuesResponse struct{}

// RemoteStartTransactionRequest is the request from server to start transaction
type RemoteStartTransactionRequest struct {
	IdTag            string            `json:"idTag"`
	ConnectorId      int               `json:"connectorId,omitempty"`
	ChargingProfile  *ChargingProfile  `json:"chargingProfile,omitempty"`
}

// ChargingProfile represents a charging profile
type ChargingProfile struct {
	ChargingProfileId      int                     `json:"chargingProfileId"`
	TransactionId          int                     `json:"transactionId,omitempty"`
	StackLevel             int                     `json:"stackLevel"`
	ChargingProfilePurpose string                  `json:"chargingProfilePurpose"`
	ChargingProfileKind    string                  `json:"chargingProfileKind"`
	RecurrencyKind         string                  `json:"recurrencyKind,omitempty"`
	ValidFrom              string                  `json:"validFrom,omitempty"`
	ValidTo                string                  `json:"validTo,omitempty"`
	ChargingSchedule       *ChargingSchedule       `json:"chargingSchedule"`
}

// ChargingSchedule represents a charging schedule
type ChargingSchedule struct {
	Duration               int                      `json:"duration,omitempty"`
	StartSchedule          string                   `json:"startSchedule,omitempty"`
	ChargingRateUnit       string                   `json:"chargingRateUnit"`
	ChargingSchedulePeriod []ChargingSchedulePeriod `json:"chargingSchedulePeriod"`
	MinChargingRate        float64                  `json:"minChargingRate,omitempty"`
}

// ChargingSchedulePeriod represents a period in charging schedule
type ChargingSchedulePeriod struct {
	StartPeriod  int     `json:"startPeriod"`
	Limit        float64 `json:"limit"`
	NumberPhases int     `json:"numberPhases,omitempty"`
}

// RemoteStartTransactionResponse is the response to RemoteStartTransaction
type RemoteStartTransactionResponse struct {
	Status string `json:"status"` // Accepted, Rejected
}

// RemoteStopTransactionRequest is the request from server to stop transaction
type RemoteStopTransactionRequest struct {
	TransactionId int `json:"transactionId"`
}

// RemoteStopTransactionResponse is the response to RemoteStopTransaction
type RemoteStopTransactionResponse struct {
	Status string `json:"status"` // Accepted, Rejected
}

// HeartbeatRequest is the request for Heartbeat
type HeartbeatRequest struct{}

// HeartbeatResponse is the response for Heartbeat
type HeartbeatResponse struct {
	CurrentTime string `json:"currentTime"`
}

// DataTransferRequest is the request for DataTransfer
type DataTransferRequest struct {
	VendorId  string `json:"vendorId"`
	MessageId string `json:"messageId,omitempty"`
	Data      string `json:"data,omitempty"`
}

// DataTransferResponse is the response for DataTransfer
type DataTransferResponse struct {
	Status string `json:"status"` // Accepted, Rejected, UnknownMessageId, UnknownVendorId
	Data   string `json:"data,omitempty"`
}

// Call represents an OCPP Call message [MessageTypeId, UniqueId, Action, Payload]
type Call struct {
	MessageTypeId int
	UniqueId      string
	Action        string
	Payload       interface{}
}

// CallResult represents an OCPP CallResult message [MessageTypeId, UniqueId, Payload]
type CallResult struct {
	MessageTypeId int
	UniqueId      string
	Payload       interface{}
}

// CallError represents an OCPP CallError message [MessageTypeId, UniqueId, ErrorCode, ErrorDescription, ErrorDetails]
type CallError struct {
	MessageTypeId    int
	UniqueId         string
	ErrorCode        string
	ErrorDescription string
	ErrorDetails     interface{}
}

// MarshalCall marshals a Call message to JSON
func MarshalCall(uniqueId, action string, payload interface{}) ([]byte, error) {
	msg := []interface{}{MessageTypeCall, uniqueId, action, payload}
	return json.Marshal(msg)
}

// MarshalCallResult marshals a CallResult message to JSON
func MarshalCallResult(uniqueId string, payload interface{}) ([]byte, error) {
	msg := []interface{}{MessageTypeCallResult, uniqueId, payload}
	return json.Marshal(msg)
}

// MarshalCallError marshals a CallError message to JSON
func MarshalCallError(uniqueId, errorCode, errorDescription string, errorDetails interface{}) ([]byte, error) {
	msg := []interface{}{MessageTypeCallError, uniqueId, errorCode, errorDescription, errorDetails}
	return json.Marshal(msg)
}

// ParseMessage parses an OCPP message and returns its type and components
func ParseMessage(data []byte) (messageType int, uniqueId string, payload json.RawMessage, action string, err error) {
	var raw []json.RawMessage
	if err = json.Unmarshal(data, &raw); err != nil {
		return 0, "", nil, "", fmt.Errorf("failed to parse message: %w", err)
	}

	if len(raw) < 3 {
		return 0, "", nil, "", fmt.Errorf("invalid message format: expected at least 3 elements")
	}

	if err = json.Unmarshal(raw[0], &messageType); err != nil {
		return 0, "", nil, "", fmt.Errorf("failed to parse message type: %w", err)
	}

	if err = json.Unmarshal(raw[1], &uniqueId); err != nil {
		return 0, "", nil, "", fmt.Errorf("failed to parse unique id: %w", err)
	}

	switch messageType {
	case MessageTypeCall:
		if len(raw) < 4 {
			return 0, "", nil, "", fmt.Errorf("invalid Call message: expected 4 elements")
		}
		if err = json.Unmarshal(raw[2], &action); err != nil {
			return 0, "", nil, "", fmt.Errorf("failed to parse action: %w", err)
		}
		payload = raw[3]
	case MessageTypeCallResult:
		payload = raw[2]
	case MessageTypeCallError:
		payload = raw[2] // Error code
	default:
		return 0, "", nil, "", fmt.Errorf("unknown message type: %d", messageType)
	}

	return messageType, uniqueId, payload, action, nil
}
