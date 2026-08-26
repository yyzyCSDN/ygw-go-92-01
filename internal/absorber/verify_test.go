package absorber_test

import (
	"testing"

	"regenbrake/internal/absorber"
	"regenbrake/internal/model"
)

func TestRecoveryUsesLatestSnapshot(t *testing.T) {
	dev := absorber.NewDevice("A", absorber.DirectDriver{})
	dev.ApplySnapshot(model.DeviceSnapshot{ID: "A", Status: model.StatusFault, Revision: 1})
	dev.ApplySnapshot(model.DeviceSnapshot{ID: "A", Status: model.StatusRecovered, Revision: 2})
	if status := dev.Status(); status != model.StatusRecovered {
		t.Fatalf("latest snapshot should win, got %s", status)
	}
	dev.ApplySnapshot(model.DeviceSnapshot{ID: "A", Status: model.StatusFault, Revision: 1})
	if status := dev.Status(); status != model.StatusRecovered {
		t.Fatalf("stale snapshot should not overwrite recovered device, got %s", status)
	}
}
