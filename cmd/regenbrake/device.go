package main

import (
	"regenbrake/internal/absorber"
	"regenbrake/internal/alarm"
)

func buildDevices(mgr *alarm.Manager) map[string]*absorber.Device {
	ids := []string{"A01", "A02", "A03"}
	devices := make(map[string]*absorber.Device, len(ids))
	for _, id := range ids {
		dev := absorber.NewDevice(id, absorber.DirectDriver{})
		devices[id] = dev
		mgr.Register(dev)
	}
	return devices
}
