package absorber_test

import (
	"sync"
	"testing"

	"regenbrake/internal/absorber"
	"regenbrake/internal/coop"
	"regenbrake/internal/model"
	"regenbrake/internal/train"
)

func TestAbsorberStateNoOverwrite(t *testing.T) {
	reg := train.NewRegistry()
	reg.Register(model.TrainState{ID: "T1", Braking: true, Power: 100})
	reg.Register(model.TrainState{ID: "T2", Braking: true, Power: 100})
	dev := absorber.NewDevice("A", absorber.DirectDriver{})
	coord := coop.NewCoordinator(reg, dev)
	const n = 40
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = coord.Allocate([]model.BrakingReport{
				{TrainID: "T1", Power: 10, Capacity: 10},
				{TrainID: "T2", Power: 10, Capacity: 10},
			})
		}()
	}
	wg.Wait()
	want := float64(n) * 20
	if got := dev.Allocated(); got != want {
		t.Fatalf("allocated %v, want %v", got, want)
	}
}
