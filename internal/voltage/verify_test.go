package voltage_test

import (
	"context"
	"testing"
	"time"

	"regenbrake/internal/absorber"
	"regenbrake/internal/inverter"
	"regenbrake/internal/voltage"
)

func TestVoltageTimeoutRecovers(t *testing.T) {
	dev := absorber.NewDevice("A", absorber.DirectDriver{})
	inv := inverter.New()
	ctrl := voltage.NewController(dev, inv, 900)
	done := make(chan error, 1)
	go func() {
		done <- ctrl.WaitAbsorbConfirm(context.Background(), 50*time.Millisecond)
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected absorb confirm timeout")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("wait did not recover after timeout")
	}
}
