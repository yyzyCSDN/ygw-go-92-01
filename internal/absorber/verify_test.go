package absorber_test

import (
	"testing"

	"regenbrake/internal/absorber"
	"regenbrake/internal/model"
	"regenbrake/internal/train"
)

func TestTrainCacheInvalidatedOnState(t *testing.T) {
	dev := absorber.NewDevice("A", absorber.DirectDriver{})
	reg := train.NewRegistry()
	reg.Register(model.TrainState{ID: "T1", Braking: true, Power: 100})
	reg.AddListener(dev)
	if err := reg.SetBraking("T1", true, 100); err != nil {
		t.Fatal(err)
	}
	if !dev.TrainBraking("T1") {
		t.Fatal("expected braking true to be cached")
	}
	if err := reg.SetBraking("T1", false, 0); err != nil {
		t.Fatal(err)
	}
	if dev.TrainBraking("T1") {
		t.Fatal("cache not invalidated after train stops braking")
	}
}
