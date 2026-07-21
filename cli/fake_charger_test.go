package cli

// fakeCharger is an in-memory implementation of the Charger interface used by
// the command tests. It holds simple state, exposes per-method programmable
// error hooks, and records calls/arguments for assertions. It performs no
// network I/O, spawns no goroutines, and uses no timers, so tests are fully
// deterministic.
type fakeCharger struct {
	connected    bool
	status       string
	soc          float64
	current      float64
	power        float64
	charging     bool
	licensePlate string

	// Programmable errors (nil = success path).
	connectErr     error
	bootErr        error
	statusNotifErr error
	setStatusErr   error
	pluginErr      error
	unplugErr      error
	startErr       error
	stopErr        error
	meterErr       error
	plateErr       error
	setSOCErr      error
	setCurrentErr  error
	setPowerErr    error

	// Call recording.
	connectCalls    int
	disconnectCalls int
	bootCalls       int
	statusNotifArg  string
	statusNotifSet  bool
	lastStartIDTag  string
	lastStopReason  string
	meterCalls      int
	lastPlate       string
}

func (f *fakeCharger) IsConnected() bool { return f.connected }

func (f *fakeCharger) Connect() error {
	f.connectCalls++
	if f.connectErr != nil {
		return f.connectErr
	}
	f.connected = true
	return nil
}

func (f *fakeCharger) Disconnect() {
	f.disconnectCalls++
	f.connected = false
	f.charging = false
}

func (f *fakeCharger) BootNotification() error {
	f.bootCalls++
	return f.bootErr
}

func (f *fakeCharger) StatusNotification(status string) error {
	f.statusNotifArg = status
	f.statusNotifSet = true
	return f.statusNotifErr
}

func (f *fakeCharger) GetStatus() string { return f.status }

func (f *fakeCharger) SetStatus(status string) error {
	if f.setStatusErr != nil {
		return f.setStatusErr
	}
	f.status = status
	return nil
}

func (f *fakeCharger) Plugin() error {
	if f.pluginErr != nil {
		return f.pluginErr
	}
	f.status = "Preparing"
	return nil
}

func (f *fakeCharger) Unplug() error {
	if f.unplugErr != nil {
		return f.unplugErr
	}
	f.status = "Available"
	f.licensePlate = ""
	return nil
}

func (f *fakeCharger) StartTransaction(idTag string) error {
	f.lastStartIDTag = idTag
	if f.startErr != nil {
		return f.startErr
	}
	f.charging = true
	return nil
}

func (f *fakeCharger) StopTransaction(reason string) error {
	f.lastStopReason = reason
	if f.stopErr != nil {
		return f.stopErr
	}
	f.charging = false
	return nil
}

func (f *fakeCharger) MeterValues() error {
	f.meterCalls++
	return f.meterErr
}

func (f *fakeCharger) SetLicensePlateAndSend(plate string) error {
	f.lastPlate = plate
	if f.plateErr != nil {
		return f.plateErr
	}
	f.licensePlate = plate
	return nil
}

func (f *fakeCharger) GetLicensePlate() string { return f.licensePlate }

func (f *fakeCharger) SetSOC(soc float64) error {
	if f.setSOCErr != nil {
		return f.setSOCErr
	}
	f.soc = soc
	return nil
}

func (f *fakeCharger) GetSOC() float64 { return f.soc }

func (f *fakeCharger) SetCurrent(current float64) error {
	if f.setCurrentErr != nil {
		return f.setCurrentErr
	}
	f.current = current
	return nil
}

func (f *fakeCharger) GetCurrent() float64 { return f.current }

func (f *fakeCharger) SetPower(power float64) error {
	if f.setPowerErr != nil {
		return f.setPowerErr
	}
	f.power = power
	return nil
}

func (f *fakeCharger) GetPower() float64 { return f.power }

func (f *fakeCharger) IsCharging() bool { return f.charging }

// Compile-time assertion that the fake satisfies the interface under test.
var _ Charger = (*fakeCharger)(nil)
