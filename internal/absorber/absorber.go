package absorber

import (
	"fmt"
	"sync"

	"regenbrake/internal/model"
)

type Device struct {
	mu           sync.Mutex
	id           string
	driver       Driver
	state        model.AbsorberState
	status       model.DeviceStatus
	allocated    float64
	trainBraking map[string]bool
	revision     int
	confirm      chan struct{}
	energy       *EnergyMeter
	ledger       *Ledger
	metrics      *Metrics
	faultReason  string
	faults       []string
	cacheVersion int
}

func NewDevice(id string, driver Driver) *Device {
	return &Device{
		id:           id,
		driver:       driver,
		state:        model.AbsorberIdle,
		status:       model.StatusOnline,
		trainBraking: make(map[string]bool),
		confirm:      make(chan struct{}, 1),
		energy:       NewEnergyMeter(),
		ledger:       NewLedger(),
		metrics:      &Metrics{},
	}
}

func (d *Device) ID() string {
	return d.id
}

func (d *Device) Engage() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.status == model.StatusFault {
		return model.ErrDeviceFault
	}
	if err := d.engageDevice(); err != nil {
		d.status = model.StatusFault
		d.faultReason = err.Error()
		d.faults = append(d.faults, err.Error())
		d.metrics.RecordFault()
		return fmt.Errorf("engage device %s: %w", d.id, err)
	}
	d.state = model.AbsorberAbsorbing
	d.status = model.StatusOnline
	d.faultReason = ""
	d.energy.RecordAbsorb(1)
	d.metrics.RecordEngage()
	d.revision++
	return nil
}

func (d *Device) engageDevice() error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		if err := d.driver.Engage(); err == nil {
			return nil
		} else {
			last = err
		}
	}
	return last
}

func (d *Device) Discharge() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.state != model.AbsorberAbsorbing {
		return fmt.Errorf("discharge requires absorbing state, got %s", d.state)
	}
	if err := d.driver.Discharge(); err != nil {
		return fmt.Errorf("discharge device %s: %w", d.id, err)
	}
	d.state = model.AbsorberDischarging
	d.energy.RecordDischarge(0.9)
	d.metrics.RecordDischarge()
	return nil
}

func (d *Device) AddCapacity(capacity float64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if capacity < 0 {
		return
	}
	d.allocated += capacity
	d.ledger.Set("_total", d.allocated)
}

func (d *Device) Allocated() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.allocated
}

func (d *Device) ResetToIdle() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = model.AbsorberIdle
	d.allocated = 0
	d.ledger.Set("_total", 0)
}

func (d *Device) State() model.AbsorberState {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.state
}

func (d *Device) Status() model.DeviceStatus {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.status
}

func (d *Device) SwitchTo(state model.AbsorberState) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.canSwitch(d.state, state) {
		return
	}
	d.state = state
}

func (d *Device) MarkFault() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.status = model.StatusFault
	d.state = model.AbsorberFault
	d.metrics.RecordFault()
}

func (d *Device) OnTrainBrakingChange(trainID string, braking bool) {
}

func (d *Device) TrainBraking(trainID string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.trainBraking[trainID]
}

func (d *Device) ConfirmSignal() <-chan struct{} {
	return d.confirm
}

func (d *Device) Confirm() {
	select {
	case d.confirm <- struct{}{}:
	default:
	}
}

func (d *Device) Energy() *EnergyMeter {
	return d.energy
}

func (d *Device) Ledger() *Ledger {
	return d.ledger
}

func (d *Device) Metrics() *Metrics {
	return d.metrics
}

func (d *Device) FaultReason() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.faultReason
}

func (d *Device) FaultCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.faults)
}

func (d *Device) CacheVersion() int {
	return 0
}

func (d *Device) canSwitch(from, to model.AbsorberState) bool {
	if from == to {
		return true
	}
	if from == model.AbsorberFault {
		return false
	}
	switch to {
	case model.AbsorberIdle:
		return from == model.AbsorberDischarging || from == model.AbsorberAbsorbing
	case model.AbsorberAbsorbing:
		return from == model.AbsorberIdle
	case model.AbsorberDischarging:
		return from == model.AbsorberAbsorbing
	default:
		return false
	}
}
