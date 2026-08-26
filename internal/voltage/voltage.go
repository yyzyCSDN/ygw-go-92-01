package voltage

import (
	"context"
	"sync"

	"regenbrake/internal/absorber"
	"regenbrake/internal/inverter"
	"regenbrake/internal/model"
)

type Controller struct {
	mu         sync.Mutex
	abs        *absorber.Device
	inv        *inverter.Inverter
	state      model.VoltageState
	voltage    float64
	threshold  float64
	hysteresis *Hysteresis
	cancelWait context.CancelFunc
	switches   int
}

func NewController(abs *absorber.Device, inv *inverter.Inverter, threshold float64) *Controller {
	return &Controller{
		abs:        abs,
		inv:        inv,
		state:      model.VoltageNormal,
		threshold:  threshold,
		hysteresis: NewHysteresis(threshold, threshold*0.95),
	}
}

func (c *Controller) State() model.VoltageState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

func (c *Controller) Voltage() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.voltage
}

func (c *Controller) Sample(voltage float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.voltage = voltage
	if c.hysteresis.Update(voltage) {
		c.state = model.VoltageHigh
	} else {
		c.state = model.VoltageNormal
	}
}

func (c *Controller) SwitchToAbsorbing() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state == model.VoltageAbsorbing {
		return nil
	}
	if err := c.abs.Engage(); err != nil {
		return err
	}
	c.syncAbsorberState(model.AbsorberAbsorbing)
	c.state = model.VoltageAbsorbing
	c.inv.Disable()
	c.switches++
	return nil
}

func (c *Controller) SwitchToRestoring() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state == model.VoltageRestoring {
		return
	}
	c.syncAbsorberState(model.AbsorberIdle)
	c.state = model.VoltageRestoring
	c.inv.Enable()
	c.switches++
}

func (c *Controller) HysteresisBand() (float64, float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hysteresis.Band()
}

func (c *Controller) High() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hysteresis.High()
}

func (c *Controller) cancelPendingWait() {
	if c.cancelWait != nil {
		c.cancelWait()
		c.cancelWait = nil
	}
}

func (c *Controller) CancelWait() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cancelPendingWait()
}

func (c *Controller) SwitchCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.switches
}

func (c *Controller) syncAbsorberState(state model.AbsorberState) {
	if c.abs.State() != state {
		c.abs.SwitchTo(state)
	}
}
