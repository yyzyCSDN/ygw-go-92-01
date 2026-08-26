package alarm_test

import (
	"errors"
	"testing"

	"regenbrake/internal/absorber"
	"regenbrake/internal/alarm"
	"regenbrake/internal/model"
)

type failingStore struct{}

func (failingStore) Writeback(id string, status model.DeviceStatus) error {
	return errors.New("writeback failed")
}

func TestRestoreWritebackErrorNotSwallowed(t *testing.T) {
	store := failingStore{}
	mgr := alarm.NewManager(store)
	dev := absorber.NewDevice("A", absorber.DirectDriver{})
	dev.MarkFault()
	mgr.Register(dev)
	if err := mgr.RestoreDevice("A"); err == nil {
		t.Fatal("expected writeback error to propagate")
	}
	if status := dev.Status(); status != model.StatusFault {
		t.Fatalf("status %s, want fault", status)
	}
}
