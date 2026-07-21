package cli

// Charger is the subset of *charger.Charger behavior the interactive commands
// depend on. Depending on an interface (rather than the concrete type) lets
// tests substitute an in-memory fake so every command runs deterministically
// without a network connection or a live OCPP server.
//
// *charger.Charger satisfies this interface unchanged; the compile-time
// assertion lives in package main where the concrete type is wired in.
type Charger interface {
	IsConnected() bool
	Connect() error
	Disconnect()
	BootNotification() error
	StatusNotification(status string) error
	GetStatus() string
	SetStatus(status string) error
	Plugin() error
	Unplug() error
	StartTransaction(idTag string) error
	StopTransaction(reason string) error
	MeterValues() error
	SetLicensePlateAndSend(plate string) error
	GetLicensePlate() string
	SetSOC(soc float64) error
	GetSOC() float64
	SetCurrent(current float64) error
	GetCurrent() float64
	SetPower(power float64) error
	GetPower() float64
	IsCharging() bool
}
