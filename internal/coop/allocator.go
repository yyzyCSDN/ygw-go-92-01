package coop

import "regenbrake/internal/model"

type Allocator struct {
	CapacityLimit float64
}

func (a Allocator) Distribute(reports []model.BrakingReport) []model.Allocation {
	out := make([]model.Allocation, 0, len(reports))
	remaining := a.CapacityLimit
	for _, report := range reports {
		if remaining <= 0 {
			break
		}
		capacity := report.Capacity
		if capacity > remaining {
			capacity = remaining
		}
		remaining -= capacity
		out = append(out, model.Allocation{TrainID: report.TrainID, Capacity: capacity, Power: report.Power})
	}
	return out
}
