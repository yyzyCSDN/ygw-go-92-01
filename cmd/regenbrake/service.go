package main

import (
	"context"
	"sync"
	"time"

	"regenbrake/internal/absorber"
	"regenbrake/internal/alarm"
	"regenbrake/internal/coop"
	"regenbrake/internal/inverter"
	"regenbrake/internal/model"
	"regenbrake/internal/record"
	"regenbrake/internal/train"
	"regenbrake/internal/voltage"
)

type Service struct {
	mu        sync.Mutex
	absorbers map[string]*absorber.Device
	inverter  *inverter.Inverter
	voltage   *voltage.Controller
	coop      *coop.Coordinator
	trains    *train.Registry
	alarm     *alarm.Manager
	journal   *record.Journal
	schedule  *train.Schedule
	procedure *alarm.Procedure
	verifier  *record.Verifier
	loop      *voltage.Loop
	threshold float64
}

func NewService(cfg Config) *Service {
	store := alarm.NewMemoryStore()
	mgr := alarm.NewManager(store)
	devices := buildDevices(mgr)
	inv := inverter.New()
	primary := devices["A01"]
	volt := voltage.NewController(primary, inv, cfg.Threshold)
	trains := train.NewRegistry()
	coord := coop.NewCoordinator(trains, primary)
	trains.AddListener(coord)
	trains.AddListener(primary)
	journal := record.NewJournal(cfg.DataDir)
	schedule := train.NewSchedule()
	procedure := alarm.NewProcedure()
	verifier := record.NewVerifier(cfg.DataDir)
	loop := voltage.NewLoop(volt, 100*time.Millisecond)
	service := &Service{
		absorbers: devices,
		inverter:  inv,
		voltage:   volt,
		coop:      coord,
		trains:    trains,
		alarm:     mgr,
		journal:   journal,
		schedule:  schedule,
		procedure: procedure,
		verifier:  verifier,
		loop:      loop,
		threshold: cfg.Threshold,
	}
	return service
}

func (s *Service) SampleVoltage(value float64) error {
	s.voltage.Sample(value)
	if s.voltage.OverThreshold() {
		if err := s.voltage.SwitchToAbsorbing(); err != nil {
			return err
		}
		return s.alarm.Engage("A01")
	}
	s.voltage.SwitchToRestoring()
	return nil
}

func (s *Service) RegisterTrain(state model.TrainState) {
	s.trains.Register(state)
}

func (s *Service) SetBraking(trainID string, braking bool, power float64) error {
	return s.trains.SetBraking(trainID, braking, power)
}

func (s *Service) Allocate(reports []model.BrakingReport) ([]model.Allocation, error) {
	allocations, err := s.coop.Allocate(reports)
	if err != nil {
		return nil, err
	}
	for _, allocation := range allocations {
		if device := s.absorbers["A01"]; device != nil {
			device.Ledger().Set(allocation.TrainID, allocation.Capacity)
		}
	}
	return allocations, nil
}

func (s *Service) RestoreDevice(deviceID string) error {
	return s.alarm.RestoreDevice(deviceID)
}

func (s *Service) Recover(snapshot model.SystemSnapshot) error {
	return s.alarm.Recover(snapshot)
}

func (s *Service) AppendRecord(entry string) error {
	return s.journal.Append(entry)
}

func (s *Service) WaitConfirm(ctx context.Context, timeout time.Duration) error {
	return s.voltage.WaitAbsorbConfirm(ctx, timeout)
}

func (s *Service) CancelAbsorbWait() {
	s.voltage.CancelWait()
}

func (s *Service) Status() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	energy := 0.0
	for _, device := range s.absorbers {
		energy += device.Energy().Absorbed()
	}
	lower, upper := s.voltage.HysteresisBand()
	stats := s.journal.Stats()
	return map[string]any{
		"voltage":    s.voltage.Voltage(),
		"state":      s.voltage.State().String(),
		"allocated":  s.coop.Total(),
		"trains":     s.trains.Count(),
		"alarms":     s.alarm.AlarmCount(),
		"inverter":   s.inverter.Active(),
		"record_seq": s.journal.Sequence(),
		"energy":     energy,
		"band_lower": lower,
		"band_upper": upper,
		"records":    stats.Records,
		"rotations":  stats.Rotations,
		"switches":   s.voltage.SwitchCount(),
	}
}

func (s *Service) Close() error {
	return s.journal.Close()
}

func (s *Service) Snapshot() model.SystemSnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	devices := make([]model.DeviceSnapshot, 0, len(s.absorbers))
	for _, device := range s.absorbers {
		devices = append(devices, device.Snapshot())
	}
	return model.SystemSnapshot{Revision: s.alarm.Revision() + 1, Devices: devices}
}

func (s *Service) EnergyReport() map[string]map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]map[string]any, len(s.absorbers))
	for id, device := range s.absorbers {
		out[id] = map[string]any{
			"absorbed":   device.Energy().Absorbed(),
			"discharged": device.Energy().Discharged(),
			"cycles":     device.Energy().Cycles(),
			"efficiency": device.Energy().Efficiency(),
		}
	}
	return out
}

func (s *Service) AddSchedule(trainID string, station int) {
	s.schedule.Add(trainID, station)
}

func (s *Service) ScheduleStation(trainID string) int {
	return s.schedule.Station(trainID)
}

func (s *Service) TrainsAt(station int) []string {
	return s.schedule.TrainsAt(station)
}

func (s *Service) ProcedureSteps() []string {
	return s.procedure.Steps()
}

func (s *Service) RestoreSnapshot(snapshot model.SystemSnapshot) error {
	return s.alarm.RestoreSnapshot(snapshot, s.procedure)
}

func (s *Service) VerifyRecords() (int, error) {
	return s.verifier.Verify()
}

func (s *Service) TickControl(value float64) map[string]any {
	return s.loop.Tick(value)
}

func (s *Service) LedgerReport() map[string]map[string]float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]map[string]float64, len(s.absorbers))
	for id, device := range s.absorbers {
		ledger := device.Ledger()
		entry := make(map[string]float64)
		for _, trainID := range ledger.Keys() {
			entry[trainID] = ledger.Get(trainID)
		}
		entry["_total"] = ledger.Total()
		out[id] = entry
	}
	return out
}

func (s *Service) ResetEnergy(deviceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	device, ok := s.absorbers[deviceID]
	if !ok {
		return false
	}
	device.Energy().Reset()
	return true
}

func (s *Service) InverterStep() float64 {
	return s.inverter.Step()
}

func (s *Service) InverterSettled() bool {
	return s.inverter.Settled()
}

func (s *Service) AlarmEvents() []alarm.Event {
	return s.alarm.History().OpenEvents()
}

func (s *Service) ConfirmAbsorb() {
	if device := s.absorbers["A01"]; device != nil {
		device.Confirm()
	}
}

func (s *Service) MarkDeviceFault(deviceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	device, ok := s.absorbers[deviceID]
	if !ok {
		return false
	}
	device.MarkFault()
	s.alarm.Escalation().Set(deviceID, alarm.LevelCritical)
	return true
}

func (s *Service) ResetDeviceToIdle(deviceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	device, ok := s.absorbers[deviceID]
	if !ok {
		return false
	}
	device.ResetToIdle()
	return true
}

func (s *Service) DeviceAllocated(deviceID string) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	device, ok := s.absorbers[deviceID]
	if !ok {
		return 0
	}
	return device.Allocated()
}

func (s *Service) TrainBrakingState(trainID string) bool {
	device := s.absorbers["A01"]
	if device == nil {
		return false
	}
	return device.TrainBraking(trainID)
}

func (s *Service) TrainBrakingOrFalse(trainID string) bool {
	return s.trains.BrakingOrFalse(trainID)
}

func (s *Service) DeviceFaultCount(deviceID string) int {
	device, ok := s.absorbers[deviceID]
	if !ok {
		return 0
	}
	return device.FaultCount()
}

func (s *Service) CoopMissCount() int {
	return s.coop.MissCount()
}

func (s *Service) JournalClosed() bool {
	return s.journal.Closed()
}

func (s *Service) DeviceFaultReason(deviceID string) string {
	device, ok := s.absorbers[deviceID]
	if !ok {
		return ""
	}
	return device.FaultReason()
}

func (s *Service) DeviceCacheVersion(deviceID string) int {
	device, ok := s.absorbers[deviceID]
	if !ok {
		return 0
	}
	return device.CacheVersion()
}

func (s *Service) RestoreCount() int {
	return s.alarm.RestoreCount()
}

func (s *Service) RotationCount() int {
	return s.journal.RotationCount()
}

func (s *Service) RotateRecord() error {
	return s.journal.Rotate()
}

func (s *Service) RecordEntryCount() int {
	return s.journal.EntryCount()
}

func (s *Service) RebalanceCoop(limit float64) []model.Allocation {
	return s.coop.Rebalance(limit)
}

func (s *Service) CoopAllocations() []model.Allocation {
	return s.coop.Allocations()
}

func (s *Service) CoopCapacityPerTrain() map[string]float64 {
	return s.coop.CapacityPerTrain()
}

func (s *Service) TrainList() []model.TrainState {
	return s.trains.List()
}

func (s *Service) BrakingTrainList() []model.TrainState {
	return s.trains.BrakingTrains()
}

func (s *Service) TrainExists(trainID string) bool {
	return s.trains.Exists(trainID)
}

func (s *Service) InverterCurrent() float64 {
	return s.inverter.Current()
}

func (s *Service) HysteresisHigh() bool {
	return s.voltage.High()
}

func (s *Service) DeviceHasAlarm(deviceID string) bool {
	return s.alarm.HasAlarm(deviceID)
}

func (s *Service) ControlInterval() time.Duration {
	return s.loop.Interval()
}

func (s *Service) DeviceMetrics(deviceID string) *absorber.Metrics {
	s.mu.Lock()
	defer s.mu.Unlock()
	device, ok := s.absorbers[deviceID]
	if !ok {
		return nil
	}
	return device.Metrics()
}

func (s *Service) RestoreControl() {
	s.voltage.Restore()
}

func (s *Service) SetInverterFeedback(power float64) {
	s.inverter.SetFeedback(power)
}

func (s *Service) FeedbackPower() float64 {
	return s.inverter.FeedbackPower()
}

func (s *Service) FeedInverter(absorbed float64) float64 {
	return s.inverter.Feed(absorbed)
}

func (s *Service) CriticalAlarms() []string {
	return s.alarm.Escalation().Critical()
}
