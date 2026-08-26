package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"regenbrake/internal/coop"
	"regenbrake/internal/model"
)

func (s *Service) energyHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.EnergyReport())
}

func (s *Service) snapshotHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Snapshot())
}

func (s *Service) simulateHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Trains []model.TrainState `json:"trains"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	sim := NewSimulator(s)
	writeJSON(w, http.StatusOK, sim.Run(body.Trains))
}

func (s *Service) recordStatsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"stats":  s.journal.Stats(),
		"closed": s.JournalClosed(),
	})
}

func (s *Service) reportHandler(w http.ResponseWriter, r *http.Request) {
	reporter := NewReporter(s)
	writeJSON(w, http.StatusOK, reporter.Report())
}

func (s *Service) estimateHandler(w http.ResponseWriter, r *http.Request) {
	trainID := r.URL.Query().Get("id")
	speed, err := strconv.ParseFloat(r.URL.Query().Get("speed"), 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid speed"})
		return
	}
	reporter := NewReporter(s)
	writeJSON(w, http.StatusOK, reporter.Estimate(trainID, speed))
}

func (s *Service) planHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Reports []model.BrakingReport `json:"reports"`
		Limit   float64               `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	reporter := NewReporter(s)
	plan := coop.Plan{Reports: body.Reports}
	allocations := reporter.Plan(body.Reports, body.Limit)
	proportional := reporter.Proportional(body.Reports, body.Limit)
	writeJSON(w, http.StatusOK, map[string]any{
		"total":        reporter.PlanTotal(plan),
		"priority":     allocations,
		"proportional": proportional,
	})
}
