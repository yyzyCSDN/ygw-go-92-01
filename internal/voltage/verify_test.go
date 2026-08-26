package voltage_test

import (
	"sync"
	"testing"

	"regenbrake/internal/absorber"
	"regenbrake/internal/inverter"
	"regenbrake/internal/model"
	"regenbrake/internal/voltage"
)

func TestVoltAbsorberStateConsistent(t *testing.T) {
	for i := 0; i < 40; i++ {
		dev := absorber.NewDevice("A", absorber.DirectDriver{})
		inv := inverter.New()
		ctrl := voltage.NewController(dev, inv, 900)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = ctrl.SwitchToAbsorbing()
		}()
		go func() {
			defer wg.Done()
			ctrl.SwitchToRestoring()
		}()
		wg.Wait()
		vs := ctrl.State()
		as := dev.State()
		if vs == model.VoltageAbsorbing && as != model.AbsorberAbsorbing {
			t.Fatalf("voltage absorbing but absorber %s", as)
		}
		if vs == model.VoltageRestoring && as != model.AbsorberIdle {
			t.Fatalf("voltage restoring but absorber %s", as)
		}
	}
}
