package voltage

import "regenbrake/internal/model"

func (c *Controller) Restore() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cancelPendingWait()
	c.abs.SwitchTo(model.AbsorberIdle)
	c.state = model.VoltageRestoring
	c.inv.Enable()
	c.switches++
}

func (c *Controller) OverThreshold() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.voltage >= c.threshold
}
