package main

import (
	"encoding/json"
	"net/http"

	"regenbrake/internal/model"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Service) statusHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Status())
}

func (s *Service) sampleVoltageHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Voltage float64 `json:"voltage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !model.ValidVoltage(body.Voltage) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid voltage"})
		return
	}
	if err := s.SampleVoltage(body.Voltage); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, s.Status())
}

func (s *Service) registerTrainHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID    string  `json:"id"`
		Power float64 `json:"power"`
		Speed float64 `json:"speed"`
		Mass  float64 `json:"mass"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !model.ValidTrainID(body.ID) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid train id"})
		return
	}
	s.RegisterTrain(model.TrainState{ID: body.ID, Power: body.Power, Speed: body.Speed, Mass: body.Mass})
	writeJSON(w, http.StatusOK, map[string]string{"status": "registered"})
}

func (s *Service) setBrakingHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID      string  `json:"id"`
		Braking bool    `json:"braking"`
		Power   float64 `json:"power"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !model.ValidTrainID(body.ID) || !model.ValidCapacity(body.Power) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}
	if err := s.SetBraking(body.ID, body.Braking, body.Power); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Service) allocateHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Reports []model.BrakingReport `json:"reports"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	allocs, err := s.Allocate(body.Reports)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"allocations": allocs})
}

func (s *Service) restoreDeviceHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("id")
	if deviceID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing id"})
		return
	}
	if err := s.RestoreDevice(deviceID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "restored"})
}

func writeJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(value)
}
