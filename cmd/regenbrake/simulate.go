package main

import (
	"regenbrake/internal/coop"
	"regenbrake/internal/model"
	"regenbrake/internal/train"
)

type Simulator struct {
	service  *Service
	profile  train.Profile
	strategy coop.Strategy
}

func NewSimulator(service *Service) *Simulator {
	return &Simulator{
		service:  service,
		profile:  train.NewDefaultProfile(),
		strategy: coop.PriorityStrategy{Limit: 600},
	}
}

func (s *Simulator) Run(trains []model.TrainState) map[string]any {
	reports := make([]model.BrakingReport, 0, len(trains))
	for _, item := range trains {
		s.service.RegisterTrain(item)
		power := s.profile.BrakingPower(item.Speed)
		report := model.BrakingReport{
			TrainID:  item.ID,
			Power:    power,
			Capacity: s.profile.RequiredCapacity(power),
		}
		reports = append(reports, report)
		if err := s.service.SetBraking(item.ID, true, power); err != nil {
			continue
		}
	}
	allocations := s.strategy.Distribute(reports, 0)
	allocated := 0.0
	for _, allocation := range allocations {
		allocated += allocation.Capacity
	}
	if _, err := s.service.Allocate(reports); err != nil {
		return map[string]any{"error": err.Error()}
	}
	if err := s.service.SampleVoltage(920.0); err != nil {
		return map[string]any{"error": err.Error()}
	}
	return map[string]any{
		"trains":    len(trains),
		"reports":   len(reports),
		"allocated": allocated,
		"status":    s.service.Status(),
	}
}
