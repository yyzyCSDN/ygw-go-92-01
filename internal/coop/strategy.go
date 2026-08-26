package coop

import "regenbrake/internal/model"

type Strategy interface {
	Distribute(reports []model.BrakingReport, limit float64) []model.Allocation
}

type ProportionalStrategy struct {
	Limit float64
}

func (s ProportionalStrategy) Distribute(reports []model.BrakingReport, limit float64) []model.Allocation {
	capacity := limit
	if capacity <= 0 {
		capacity = s.Limit
	}
	total := 0.0
	for _, report := range reports {
		total += report.Capacity
	}
	if total <= 0 {
		return nil
	}
	out := make([]model.Allocation, 0, len(reports))
	for _, report := range reports {
		ratio := report.Capacity / total
		share := capacity * ratio
		out = append(out, model.Allocation{TrainID: report.TrainID, Capacity: share, Power: report.Power})
	}
	return out
}

type PriorityStrategy struct {
	Limit float64
}

func (s PriorityStrategy) Distribute(reports []model.BrakingReport, limit float64) []model.Allocation {
	capacity := limit
	if capacity <= 0 {
		capacity = s.Limit
	}
	out := make([]model.Allocation, 0, len(reports))
	remaining := capacity
	for _, report := range reports {
		if remaining <= 0 {
			break
		}
		share := report.Capacity
		if share > remaining {
			share = remaining
		}
		remaining -= share
		out = append(out, model.Allocation{TrainID: report.TrainID, Capacity: share, Power: report.Power})
	}
	return out
}
