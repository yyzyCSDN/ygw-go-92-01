package absorber

import (
	"fmt"

	"regenbrake/internal/model"
)

func (d *Device) ApplySnapshot(snap model.DeviceSnapshot) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if snap.Revision < d.revision {
		return
	}
	d.status = snap.Status
	d.state = snap.State
	d.allocated = snap.Capacity
	d.revision = snap.Revision
}

func (d *Device) Revision() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.revision
}

func (d *Device) Restore() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.status == model.StatusFault {
		return model.ErrDeviceFault
	}
	d.status = model.StatusRecovered
	d.state = model.AbsorberIdle
	d.metrics.RecordRecover()
	return nil
}

func (d *Device) Snapshot() model.DeviceSnapshot {
	d.mu.Lock()
	defer d.mu.Unlock()
	return model.DeviceSnapshot{
		ID:       d.id,
		Status:   d.status,
		State:    d.state,
		Capacity: d.allocated,
		Revision: d.revision,
	}
}

func (d *Device) Describe() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return fmt.Sprintf("%s:%s:%s", d.id, d.status, d.state)
}

type CapacitySnapshot struct {
	Allocated float64
	State     model.AbsorberState
	Status    model.DeviceStatus
}

func (d *Device) CapacitySnapshot() CapacitySnapshot {
	d.mu.Lock()
	defer d.mu.Unlock()
	return CapacitySnapshot{
		Allocated: d.allocated,
		State:     d.state,
		Status:    d.status,
	}
}
