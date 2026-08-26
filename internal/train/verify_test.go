package train_test

import (
	"testing"

	"regenbrake/internal/absorber"
	"regenbrake/internal/coop"
	"regenbrake/internal/model"
	"regenbrake/internal/train"
)

func TestEmptyTrainNoNilPanic(t *testing.T) {
	reg := train.NewRegistry()
	dev := absorber.NewDevice("A", absorber.DirectDriver{})
	coord := coop.NewCoordinator(reg, dev)
	allocations, err := coord.Allocate([]model.BrakingReport{{TrainID: "missing", Power: 10, Capacity: 10}})
	if err != nil {
		t.Fatalf("allocate missing train: %v", err)
	}
	if len(allocations) != 0 {
		t.Fatalf("expected no allocations for missing train, got %d", len(allocations))
	}
}
