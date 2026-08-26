package main

import (
	"regenbrake/internal/coop"
	"regenbrake/internal/model"
	"regenbrake/internal/train"
)

type Reporter struct {
	service *Service
	profile train.Profile
}

func NewReporter(service *Service) *Reporter {
	return &Reporter{service: service, profile: train.NewDefaultProfile()}
}

func (r *Reporter) Report() map[string]any {
	return map[string]any{
		"status": r.service.Status(),
		"energy": r.service.EnergyReport(),
		"ledger": r.service.LedgerReport(),
		"alarms": r.service.AlarmEvents(),
	}
}

func (r *Reporter) Estimate(trainID string, speed float64) map[string]float64 {
	return r.profile.Describe(model.TrainState{ID: trainID, Speed: speed})
}

func (r *Reporter) Plan(reports []model.BrakingReport, limit float64) []model.Allocation {
	allocator := coop.Allocator{CapacityLimit: limit}
	return allocator.Distribute(reports)
}

func (r *Reporter) PlanTotal(plan coop.Plan) float64 {
	return plan.TotalCapacity()
}

func (r *Reporter) Proportional(reports []model.BrakingReport, limit float64) []model.Allocation {
	strategy := coop.ProportionalStrategy{Limit: limit}
	return strategy.Distribute(reports, 0)
}
