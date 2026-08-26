package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"regenbrake/internal/model"
)

func (s *Service) scheduleHandler(w http.ResponseWriter, r *http.Request) {
	trainID := r.URL.Query().Get("id")
	station, err := strconv.Atoi(r.URL.Query().Get("station"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid station"})
		return
	}
	s.AddSchedule(trainID, station)
	writeJSON(w, http.StatusOK, map[string]any{"id": trainID, "station": station})
}

func (s *Service) scheduleListHandler(w http.ResponseWriter, r *http.Request) {
	station, err := strconv.Atoi(r.URL.Query().Get("station"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid station"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"trains": s.TrainsAt(station)})
}

func (s *Service) procedureHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"steps": s.ProcedureSteps()})
}

func (s *Service) restoreSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Revision int                    `json:"revision"`
		Devices  []model.DeviceSnapshot `json:"devices"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	snapshot := model.SystemSnapshot{Revision: body.Revision, Devices: body.Devices}
	if err := s.RestoreSnapshot(snapshot); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "restored"})
}

func (s *Service) verifyRecordsHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := s.VerifyRecords()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"bytes": bytes})
}

func (s *Service) tickHandler(w http.ResponseWriter, r *http.Request) {
	value, err := strconv.ParseFloat(r.URL.Query().Get("voltage"), 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid voltage"})
		return
	}
	writeJSON(w, http.StatusOK, s.TickControl(value))
}

func (s *Service) ledgerHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.LedgerReport())
}

func (s *Service) resetEnergyHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("id")
	if !s.ResetEnergy(deviceID) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "device not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "reset"})
}

func (s *Service) inverterStepHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"step":    s.InverterStep(),
		"current": s.InverterStep(),
		"settled": s.InverterSettled(),
	})
}

func (s *Service) alarmsHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("id")
	byDevice := s.alarm.History().ByDevice(deviceID)
	hasAlarm := s.DeviceHasAlarm(deviceID)
	critical := s.CriticalAlarms()
	writeJSON(w, http.StatusOK, map[string]any{
		"open":      s.AlarmEvents(),
		"by_device": byDevice,
		"has_alarm": hasAlarm,
		"critical":  critical,
	})
}

func (s *Service) confirmHandler(w http.ResponseWriter, r *http.Request) {
	s.ConfirmAbsorb()
	writeJSON(w, http.StatusOK, map[string]string{"status": "confirmed"})
}

func (s *Service) faultHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("id")
	if !s.MarkDeviceFault(deviceID) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "device not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "fault"})
}

func (s *Service) idleHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("id")
	if !s.ResetDeviceToIdle(deviceID) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "device not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "idle"})
}

func (s *Service) deviceAllocatedHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("id")
	writeJSON(w, http.StatusOK, map[string]any{"allocated": s.DeviceAllocated(deviceID)})
}

func (s *Service) rotateHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.RotateRecord(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "rotated"})
}

func (s *Service) rebalanceHandler(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.ParseFloat(r.URL.Query().Get("limit"), 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid limit"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"allocations": s.RebalanceCoop(limit)})
}

func (s *Service) coopAllocationsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"allocations": s.CoopAllocations(),
		"per_train":   s.CoopCapacityPerTrain(),
	})
}

func (s *Service) trainListHandler(w http.ResponseWriter, r *http.Request) {
	trainID := r.URL.Query().Get("id")
	writeJSON(w, http.StatusOK, map[string]any{
		"all":     s.TrainList(),
		"braking": s.BrakingTrainList(),
		"exists":  s.TrainExists(trainID),
	})
}

func (s *Service) controlInfoHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"interval":   s.ControlInterval().String(),
		"hysteresis": s.HysteresisHigh(),
		"inverter":   s.InverterCurrent(),
		"records":    s.RecordEntryCount(),
		"switches":   s.voltage.SwitchCount(),
		"cache_ver":  s.DeviceCacheVersion("A01"),
		"restores":   s.RestoreCount(),
		"rotations":  s.RotationCount(),
	})
}

func (s *Service) metricsHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("id")
	metrics := s.DeviceMetrics(deviceID)
	if metrics == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "device not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"engaged":      metrics.Engaged,
		"discharged":   metrics.Discharged,
		"recovered":    metrics.Recovered,
		"faults":       metrics.Faults,
		"fault_count":  s.DeviceFaultCount(deviceID),
		"fault_reason": s.DeviceFaultReason(deviceID),
	})
}

func (s *Service) coopStatsHandler(w http.ResponseWriter, r *http.Request) {
	trainID := r.URL.Query().Get("id")
	writeJSON(w, http.StatusOK, map[string]any{
		"misses":  s.CoopMissCount(),
		"stale":   s.coop.StaleReleaseCount(),
		"missing": s.coop.MissingTrains(),
		"braking": s.TrainBrakingOrFalse(trainID),
	})
}

func (s *Service) restoreControlHandler(w http.ResponseWriter, r *http.Request) {
	s.RestoreControl()
	writeJSON(w, http.StatusOK, map[string]string{"status": "restored"})
}

func (s *Service) cancelWaitHandler(w http.ResponseWriter, r *http.Request) {
	s.CancelAbsorbWait()
	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (s *Service) scheduleStationHandler(w http.ResponseWriter, r *http.Request) {
	trainID := r.URL.Query().Get("id")
	writeJSON(w, http.StatusOK, map[string]any{"station": s.ScheduleStation(trainID)})
}

func (s *Service) trainBrakingHandler(w http.ResponseWriter, r *http.Request) {
	trainID := r.URL.Query().Get("id")
	writeJSON(w, http.StatusOK, map[string]any{"braking": s.TrainBrakingState(trainID)})
}

func (s *Service) inverterFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	power, err := strconv.ParseFloat(r.URL.Query().Get("power"), 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid power"})
		return
	}
	s.SetInverterFeedback(power)
	writeJSON(w, http.StatusOK, map[string]any{
		"feedback": s.FeedbackPower(),
		"feed":     s.FeedInverter(power),
	})
}

func (s *Service) lastRecordHandler(w http.ResponseWriter, r *http.Request) {
	entry, ok := s.journal.LastEntry()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no records"})
		return
	}
	writeJSON(w, http.StatusOK, entry)
}
