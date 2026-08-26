package coop

import (
	"testing"

	"regenbrake/internal/absorber"
	"regenbrake/internal/model"
	"regenbrake/internal/train"
)

// 复现原 bug：协同分配一个未登记列车的制动状态不应 panic，
// 而是返回空结果，并把缺失列车记入 MissingTrains/MissCount。
func TestAllocateUnregisteredTrainDoesNotCrash(t *testing.T) {
	trains := train.NewRegistry()
	abs := absorber.NewDevice("A01", absorber.DirectDriver{})
	coord := NewCoordinator(trains, abs)

	reports := []model.BrakingReport{
		{TrainID: "GHOST-1", Power: 500, Capacity: 400},
		{TrainID: "GHOST-2", Power: 300, Capacity: 240},
	}

	allocs, err := coord.Allocate(reports)
	if err != nil {
		t.Fatalf("未登记列车不应返回错误，got %v", err)
	}
	if len(allocs) != 0 {
		t.Fatalf("未登记列车应给空分配，got %v", allocs)
	}
	if got := coord.MissCount(); got != 2 {
		t.Fatalf("MissCount 应为 2，got %d", got)
	}
	if missing := coord.MissingTrains(); len(missing) != 2 {
		t.Fatalf("MissingTrains 应有 2 个，got %v", missing)
	}
}

// 已登记且处于制动状态的列车应正常分配。
func TestAllocateRegisteredBrakingTrain(t *testing.T) {
	trains := train.NewRegistry()
	abs := absorber.NewDevice("A01", absorber.DirectDriver{})
	coord := NewCoordinator(trains, abs)

	trains.Register(model.TrainState{ID: "T-1", Braking: true, Power: 500})

	reports := []model.BrakingReport{{TrainID: "T-1", Power: 500, Capacity: 400}}
	allocs, err := coord.Allocate(reports)
	if err != nil {
		t.Fatalf("已登记列车不应返回错误，got %v", err)
	}
	if len(allocs) != 1 || allocs[0].TrainID != "T-1" || allocs[0].Capacity != 400 {
		t.Fatalf("应分配给 T-1 400，got %v", allocs)
	}
	if got := coord.Total(); got != 400 {
		t.Fatalf("Total 应为 400，got %v", got)
	}
}

// 未登记列车与已登记制动列车混合时，只分配给后者。
func TestAllocateMixedRegisteredAndUnregistered(t *testing.T) {
	trains := train.NewRegistry()
	abs := absorber.NewDevice("A01", absorber.DirectDriver{})
	coord := NewCoordinator(trains, abs)

	trains.Register(model.TrainState{ID: "T-1", Braking: true, Power: 500})

	reports := []model.BrakingReport{
		{TrainID: "GHOST-1", Power: 999, Capacity: 999},
		{TrainID: "T-1", Power: 500, Capacity: 400},
	}
	allocs, err := coord.Allocate(reports)
	if err != nil {
		t.Fatalf("got %v", err)
	}
	if len(allocs) != 1 || allocs[0].TrainID != "T-1" {
		t.Fatalf("只应分配给 T-1，got %v", allocs)
	}
	if got := coord.MissCount(); got != 1 {
		t.Fatalf("MissCount 应为 1，got %d", got)
	}
}

// Lookup 对未登记列车应返回 (nil, false)。
func TestRegistryLookupUnregisteredReturnsFalse(t *testing.T) {
	trains := train.NewRegistry()
	if _, ok := trains.Lookup("NOPE"); ok {
		t.Fatal("未登记列车 Lookup 应返回 false")
	}
	if trains.Exists("NOPE") {
		t.Fatal("未登记列车 Exists 应返回 false")
	}
}
