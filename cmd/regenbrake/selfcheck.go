package main

import (
	"context"
	"time"

	"regenbrake/internal/model"
)

func (s *Service) SelfCheck() error {
	s.RegisterTrain(model.TrainState{ID: "T-SELF", Braking: false, Power: 0})
	if err := s.SampleVoltage(850.0); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := s.WaitConfirm(ctx, 100*time.Millisecond)
	if err != nil && err != model.ErrAbsorbTimeout {
		return err
	}
	return s.AppendRecord("selfcheck ok")
}
