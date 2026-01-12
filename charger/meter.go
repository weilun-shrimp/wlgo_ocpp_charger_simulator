package charger

import (
	"fmt"
	"log"
	"time"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v16"
	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/ocpp/v201"
)

// MeterValues updates meter values locally and sends to server if connected
func (c *Charger) MeterValues() error {
	c.mu.Lock()
	// Simulate energy consumption
	energyWh := int(c.config.MaxPower * float64(c.config.MeterValuesInterval) / 3600 * 0.5)
	c.meterValue += energyWh

	// Update SOC
	socIncrease := (float64(energyWh) / c.config.BatteryCapacity) * 100
	c.soc += socIncrease
	if c.soc > 100 {
		c.soc = 100
	}

	meterValue := c.meterValue
	soc := c.soc
	transactionId := c.transactionId
	transactionIdStr := c.transactionIdStr
	isConnected := c.isConnected
	c.seqNo++
	seqNo := c.seqNo
	c.mu.Unlock()

	log.Printf("MeterValues: energy=%d Wh, SoC=%.1f%%", meterValue, soc)

	// Send to server if connected
	if isConnected {
		if c.config.IsOCPP16() {
			return c.sendMeterValuesV16(meterValue, soc, transactionId)
		}
		return c.sendMeterValuesV201(meterValue, soc, transactionIdStr, seqNo)
	}
	return nil
}

func (c *Charger) sendMeterValuesV16(meterValue int, soc float64, transactionId int) error {
	req := v16.MeterValuesRequest{
		ConnectorId:   c.config.ConnectorID,
		TransactionId: transactionId,
		MeterValue: []v16.MeterValueEntry{
			{
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				SampledValue: []v16.SampledValue{
					{
						Value:     fmt.Sprintf("%d", meterValue),
						Context:   "Sample.Periodic",
						Measurand: "Energy.Active.Import.Register",
						Unit:      "Wh",
					},
					{
						Value:     fmt.Sprintf("%.1f", c.config.MaxCurrent*0.5),
						Context:   "Sample.Periodic",
						Measurand: "Current.Import",
						Unit:      "A",
					},
					{
						Value:     fmt.Sprintf("%.1f", c.config.MaxPower*0.5),
						Context:   "Sample.Periodic",
						Measurand: "Power.Active.Import",
						Unit:      "W",
					},
					{
						Value:     fmt.Sprintf("%.1f", soc),
						Context:   "Sample.Periodic",
						Measurand: "SoC",
						Unit:      "Percent",
					},
				},
			},
		},
	}

	_, err := c.sendCall(v16.ActionMeterValues, req)
	if err != nil {
		return fmt.Errorf("MeterValues failed: %w", err)
	}

	log.Printf("MeterValues sent: energy=%d Wh, SoC=%.1f%%", meterValue, soc)
	return nil
}

func (c *Charger) sendMeterValuesV201(meterValue int, soc float64, transactionIdStr string, seqNo int) error {
	req := v201.TransactionEventRequest{
		EventType:     v201.TransactionEventUpdated,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		TriggerReason: v201.TriggerReasonMeterValuePeriodic,
		SeqNo:         seqNo,
		TransactionInfo: v201.Transaction{
			TransactionId: transactionIdStr,
			ChargingState: v201.ChargingStateCharging,
		},
		MeterValue: []v201.MeterValue{
			{
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				SampledValue: []v201.SampledValue{
					{
						Value:     float64(meterValue),
						Context:   "Sample.Periodic",
						Measurand: "Energy.Active.Import.Register",
						UnitOfMeasure: &v201.UnitOfMeasure{
							Unit: "Wh",
						},
					},
					{
						Value:     c.config.MaxCurrent * 0.5,
						Context:   "Sample.Periodic",
						Measurand: "Current.Import",
						UnitOfMeasure: &v201.UnitOfMeasure{
							Unit: "A",
						},
					},
					{
						Value:     c.config.MaxPower * 0.5,
						Context:   "Sample.Periodic",
						Measurand: "Power.Active.Import",
						UnitOfMeasure: &v201.UnitOfMeasure{
							Unit: "W",
						},
					},
					{
						Value:     soc,
						Context:   "Sample.Periodic",
						Measurand: "SoC",
						UnitOfMeasure: &v201.UnitOfMeasure{
							Unit: "Percent",
						},
					},
				},
			},
		},
	}

	_, err := c.sendCall(v201.ActionTransactionEvent, req)
	if err != nil {
		return fmt.Errorf("TransactionEvent (Updated) failed: %w", err)
	}

	log.Printf("TransactionEvent (Updated) sent: energy=%d Wh, SoC=%.1f%%", meterValue, soc)
	return nil
}

// StartMeterValuesLoop starts auto meter updates while charging
func (c *Charger) StartMeterValuesLoop() {
	interval := time.Duration(c.config.MeterValuesInterval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	c.mu.RLock()
	stopCh := c.meterStopCh
	c.mu.RUnlock()

	log.Printf("Meter loop started")

	for {
		select {
		case <-stopCh:
			log.Printf("Meter loop stopped")
			return
		case <-ticker.C:
			if err := c.MeterValues(); err != nil {
				log.Printf("MeterValues error: %v", err)
			}
		}
	}
}
