package v201

import (
	"encoding/json"
	"fmt"
)

// OCPP 2.0.1 Message Types
const (
	MessageTypeCall       = 2 // Request
	MessageTypeCallResult = 3 // Response
	MessageTypeCallError  = 4 // Error
)

// OCPP 2.0.1 Actions
const (
	ActionBootNotification        = "BootNotification"
	ActionStatusNotification      = "StatusNotification"
	ActionTransactionEvent        = "TransactionEvent"
	ActionMeterValues             = "MeterValues"
	ActionRequestStartTransaction = "RequestStartTransaction"
	ActionRequestStopTransaction  = "RequestStopTransaction"
	ActionHeartbeat               = "Heartbeat"
	ActionDataTransfer            = "DataTransfer"
)

// ConnectorStatus represents the status of a connector in OCPP 2.0.1
type ConnectorStatus string

const (
	ConnectorStatusAvailable   ConnectorStatus = "Available"
	ConnectorStatusOccupied    ConnectorStatus = "Occupied"
	ConnectorStatusReserved    ConnectorStatus = "Reserved"
	ConnectorStatusUnavailable ConnectorStatus = "Unavailable"
	ConnectorStatusFaulted     ConnectorStatus = "Faulted"
)

// RegistrationStatus represents the registration status
type RegistrationStatus string

const (
	RegistrationAccepted RegistrationStatus = "Accepted"
	RegistrationPending  RegistrationStatus = "Pending"
	RegistrationRejected RegistrationStatus = "Rejected"
)

// TransactionEventType represents the type of transaction event
type TransactionEventType string

const (
	TransactionEventStarted TransactionEventType = "Started"
	TransactionEventUpdated TransactionEventType = "Updated"
	TransactionEventEnded   TransactionEventType = "Ended"
)

// TriggerReason represents the reason for a transaction event
type TriggerReason string

const (
	TriggerReasonAuthorized         TriggerReason = "Authorized"
	TriggerReasonCablePluggedIn     TriggerReason = "CablePluggedIn"
	TriggerReasonChargingRateChanged TriggerReason = "ChargingRateChanged"
	TriggerReasonChargingStateChanged TriggerReason = "ChargingStateChanged"
	TriggerReasonDeauthorized       TriggerReason = "Deauthorized"
	TriggerReasonEnergyLimitReached TriggerReason = "EnergyLimitReached"
	TriggerReasonEVCommunicationLost TriggerReason = "EVCommunicationLost"
	TriggerReasonEVConnectTimeout   TriggerReason = "EVConnectTimeout"
	TriggerReasonMeterValueClock    TriggerReason = "MeterValueClock"
	TriggerReasonMeterValuePeriodic TriggerReason = "MeterValuePeriodic"
	TriggerReasonTimeLimitReached   TriggerReason = "TimeLimitReached"
	TriggerReasonTrigger            TriggerReason = "Trigger"
	TriggerReasonUnlockCommand      TriggerReason = "UnlockCommand"
	TriggerReasonStopAuthorized     TriggerReason = "StopAuthorized"
	TriggerReasonEVDeparted         TriggerReason = "EVDeparted"
	TriggerReasonEVDetected         TriggerReason = "EVDetected"
	TriggerReasonRemoteStart        TriggerReason = "RemoteStart"
	TriggerReasonRemoteStop         TriggerReason = "RemoteStop"
	TriggerReasonAbnormalCondition  TriggerReason = "AbnormalCondition"
	TriggerReasonSignedDataReceived TriggerReason = "SignedDataReceived"
	TriggerReasonResetCommand       TriggerReason = "ResetCommand"
)

// ChargingState represents the charging state
type ChargingState string

const (
	ChargingStateCharging     ChargingState = "Charging"
	ChargingStateEVConnected  ChargingState = "EVConnected"
	ChargingStateSuspendedEV  ChargingState = "SuspendedEV"
	ChargingStateSuspendedEVSE ChargingState = "SuspendedEVSE"
	ChargingStateIdle         ChargingState = "Idle"
)

// ChargingStation represents a charging station
type ChargingStation struct {
	SerialNumber    string        `json:"serialNumber,omitempty"`
	Model           string        `json:"model"`
	VendorName      string        `json:"vendorName"`
	FirmwareVersion string        `json:"firmwareVersion,omitempty"`
	Modem           *Modem        `json:"modem,omitempty"`
}

// Modem represents modem information
type Modem struct {
	Iccid string `json:"iccid,omitempty"`
	Imsi  string `json:"imsi,omitempty"`
}

// BootNotificationRequest is the request for BootNotification
type BootNotificationRequest struct {
	Reason          string          `json:"reason"`
	ChargingStation ChargingStation `json:"chargingStation"`
}

// BootNotificationResponse is the response for BootNotification
type BootNotificationResponse struct {
	CurrentTime string             `json:"currentTime"`
	Interval    int                `json:"interval"`
	Status      RegistrationStatus `json:"status"`
	StatusInfo  *StatusInfo        `json:"statusInfo,omitempty"`
}

// StatusInfo provides additional status information
type StatusInfo struct {
	ReasonCode     string `json:"reasonCode"`
	AdditionalInfo string `json:"additionalInfo,omitempty"`
}

// StatusNotificationRequest is the request for StatusNotification
type StatusNotificationRequest struct {
	Timestamp       string          `json:"timestamp"`
	ConnectorStatus ConnectorStatus `json:"connectorStatus"`
	EvseId          int             `json:"evseId"`
	ConnectorId     int             `json:"connectorId"`
}

// StatusNotificationResponse is the response for StatusNotification
type StatusNotificationResponse struct{}

// EVSE represents an EVSE
type EVSE struct {
	Id          int `json:"id"`
	ConnectorId int `json:"connectorId,omitempty"`
}

// Transaction represents a transaction
type Transaction struct {
	TransactionId     string         `json:"transactionId"`
	ChargingState     ChargingState  `json:"chargingState,omitempty"`
	TimeSpentCharging int            `json:"timeSpentCharging,omitempty"`
	StoppedReason     string         `json:"stoppedReason,omitempty"`
	RemoteStartId     int            `json:"remoteStartId,omitempty"`
}

// IdToken represents an ID token
type IdToken struct {
	IdToken string `json:"idToken"`
	Type    string `json:"type"`
}

// TransactionEventRequest is the request for TransactionEvent
type TransactionEventRequest struct {
	EventType         TransactionEventType `json:"eventType"`
	Timestamp         string               `json:"timestamp"`
	TriggerReason     TriggerReason        `json:"triggerReason"`
	SeqNo             int                  `json:"seqNo"`
	TransactionInfo   Transaction          `json:"transactionInfo"`
	Offline           bool                 `json:"offline,omitempty"`
	NumberOfPhasesUsed int                 `json:"numberOfPhasesUsed,omitempty"`
	CableMaxCurrent   float64              `json:"cableMaxCurrent,omitempty"`
	ReservationId     int                  `json:"reservationId,omitempty"`
	Evse              *EVSE                `json:"evse,omitempty"`
	IdToken           *IdToken             `json:"idToken,omitempty"`
	MeterValue        []MeterValue         `json:"meterValue,omitempty"`
}

// TransactionEventResponse is the response for TransactionEvent
type TransactionEventResponse struct {
	TotalCost                 float64          `json:"totalCost,omitempty"`
	ChargingPriority          int              `json:"chargingPriority,omitempty"`
	IdTokenInfo               *IdTokenInfo     `json:"idTokenInfo,omitempty"`
	UpdatedPersonalMessage    *MessageContent  `json:"updatedPersonalMessage,omitempty"`
}

// IdTokenInfo contains authorization information
type IdTokenInfo struct {
	Status              string   `json:"status"`
	CacheExpiryDateTime string   `json:"cacheExpiryDateTime,omitempty"`
	ChargingPriority    int      `json:"chargingPriority,omitempty"`
	Language1           string   `json:"language1,omitempty"`
	Language2           string   `json:"language2,omitempty"`
	GroupIdToken        *IdToken `json:"groupIdToken,omitempty"`
	PersonalMessage     *MessageContent `json:"personalMessage,omitempty"`
}

// MessageContent represents a message
type MessageContent struct {
	Format   string `json:"format"`
	Language string `json:"language,omitempty"`
	Content  string `json:"content"`
}

// MeterValue represents meter values
type MeterValue struct {
	Timestamp    string         `json:"timestamp"`
	SampledValue []SampledValue `json:"sampledValue"`
}

// SampledValue represents a sampled value
type SampledValue struct {
	Value              float64          `json:"value"`
	Context            string           `json:"context,omitempty"`
	Measurand          string           `json:"measurand,omitempty"`
	Phase              string           `json:"phase,omitempty"`
	Location           string           `json:"location,omitempty"`
	SignedMeterValue   *SignedMeterValue `json:"signedMeterValue,omitempty"`
	UnitOfMeasure      *UnitOfMeasure   `json:"unitOfMeasure,omitempty"`
}

// SignedMeterValue represents a signed meter value
type SignedMeterValue struct {
	SignedMeterData  string `json:"signedMeterData"`
	SigningMethod    string `json:"signingMethod"`
	EncodingMethod   string `json:"encodingMethod"`
	PublicKey        string `json:"publicKey"`
}

// UnitOfMeasure represents a unit of measure
type UnitOfMeasure struct {
	Unit      string `json:"unit,omitempty"`
	Multiplier int   `json:"multiplier,omitempty"`
}

// MeterValuesRequest is the request for MeterValues
type MeterValuesRequest struct {
	EvseId     int          `json:"evseId"`
	MeterValue []MeterValue `json:"meterValue"`
}

// MeterValuesResponse is the response for MeterValues
type MeterValuesResponse struct{}

// RequestStartTransactionRequest is the request from server to start transaction
type RequestStartTransactionRequest struct {
	IdToken           IdToken           `json:"idToken"`
	RemoteStartId     int               `json:"remoteStartId"`
	EvseId            int               `json:"evseId,omitempty"`
	GroupIdToken      *IdToken          `json:"groupIdToken,omitempty"`
	ChargingProfile   *ChargingProfile  `json:"chargingProfile,omitempty"`
}

// ChargingProfile represents a charging profile
type ChargingProfile struct {
	Id                     int                     `json:"id"`
	StackLevel             int                     `json:"stackLevel"`
	ChargingProfilePurpose string                  `json:"chargingProfilePurpose"`
	ChargingProfileKind    string                  `json:"chargingProfileKind"`
	RecurrencyKind         string                  `json:"recurrencyKind,omitempty"`
	ValidFrom              string                  `json:"validFrom,omitempty"`
	ValidTo                string                  `json:"validTo,omitempty"`
	ChargingSchedule       []ChargingSchedule      `json:"chargingSchedule"`
	TransactionId          string                  `json:"transactionId,omitempty"`
}

// ChargingSchedule represents a charging schedule
type ChargingSchedule struct {
	Id                     int                      `json:"id"`
	StartSchedule          string                   `json:"startSchedule,omitempty"`
	Duration               int                      `json:"duration,omitempty"`
	ChargingRateUnit       string                   `json:"chargingRateUnit"`
	ChargingSchedulePeriod []ChargingSchedulePeriod `json:"chargingSchedulePeriod"`
	MinChargingRate        float64                  `json:"minChargingRate,omitempty"`
	SalesTariff            *SalesTariff             `json:"salesTariff,omitempty"`
}

// ChargingSchedulePeriod represents a period in charging schedule
type ChargingSchedulePeriod struct {
	StartPeriod       int     `json:"startPeriod"`
	Limit             float64 `json:"limit"`
	NumberPhases      int     `json:"numberPhases,omitempty"`
	PhaseToUse        int     `json:"phaseToUse,omitempty"`
}

// SalesTariff represents sales tariff information
type SalesTariff struct {
	Id                     int                 `json:"id"`
	SalesTariffDescription string              `json:"salesTariffDescription,omitempty"`
	NumEPriceLevels        int                 `json:"numEPriceLevels,omitempty"`
	SalesTariffEntry       []SalesTariffEntry  `json:"salesTariffEntry"`
}

// SalesTariffEntry represents a sales tariff entry
type SalesTariffEntry struct {
	RelativeTimeInterval RelativeTimeInterval `json:"relativeTimeInterval"`
	EPriceLevel          int                  `json:"ePriceLevel,omitempty"`
	ConsumptionCost      []ConsumptionCost    `json:"consumptionCost,omitempty"`
}

// RelativeTimeInterval represents a relative time interval
type RelativeTimeInterval struct {
	Start    int `json:"start"`
	Duration int `json:"duration,omitempty"`
}

// ConsumptionCost represents consumption cost
type ConsumptionCost struct {
	StartValue float64 `json:"startValue"`
	Cost       []Cost  `json:"cost"`
}

// Cost represents a cost
type Cost struct {
	CostKind       string `json:"costKind"`
	Amount         int    `json:"amount"`
	AmountMultiplier int  `json:"amountMultiplier,omitempty"`
}

// RequestStartTransactionResponse is the response to RequestStartTransaction
type RequestStartTransactionResponse struct {
	Status        string      `json:"status"` // Accepted, Rejected
	TransactionId string      `json:"transactionId,omitempty"`
	StatusInfo    *StatusInfo `json:"statusInfo,omitempty"`
}

// RequestStopTransactionRequest is the request from server to stop transaction
type RequestStopTransactionRequest struct {
	TransactionId string `json:"transactionId"`
}

// RequestStopTransactionResponse is the response to RequestStopTransaction
type RequestStopTransactionResponse struct {
	Status     string      `json:"status"` // Accepted, Rejected
	StatusInfo *StatusInfo `json:"statusInfo,omitempty"`
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
		payload = raw[2]
	default:
		return 0, "", nil, "", fmt.Errorf("unknown message type: %d", messageType)
	}

	return messageType, uniqueId, payload, action, nil
}
