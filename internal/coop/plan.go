package coop

import "regenbrake/internal/model"

type Plan struct {
	Reports []model.BrakingReport
}

func (p Plan) TotalCapacity() float64 {
	total := 0.0
	for _, report := range p.Reports {
		total += report.Capacity
	}
	return total
}
