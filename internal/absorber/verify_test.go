package absorber_test

import (
	"testing"

	"regenbrake/internal/absorber"
	"regenbrake/internal/alarm"
	"regenbrake/internal/model"
)

func TestAbsorbErrorNotSwallowed(t *testing.T) {
	store := alarm.NewMemoryStore()
	mgr := alarm.NewManager(store)
	dev := absorber.NewDevice("A", absorber.FaultyDriver{})
	mgr.Register(dev)
	if err := mgr.Engage("A"); err == nil {
		t.Fatal("expected engage error to propagate")
	}
	if !mgr.HasAlarm("A") {
		t.Fatal("expected alarm raised on engage failure")
	}
	if state := dev.State(); state != model.AbsorberIdle {
		t.Fatalf("state %s, want idle", state)
	}
}
