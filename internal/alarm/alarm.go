package alarm

import (
	"fmt"
	"sync"

	"regenbrake/internal/absorber"
	"regenbrake/internal/model"
)

type StatusStore interface {
	Writeback(id string, status model.DeviceStatus) error
}

type Manager struct {
	mu         sync.Mutex
	store      StatusStore
	absorbers  map[string]*absorber.Device
	alarms     map[string]string
	revision   int
	history    *History
	escalation *Escalation
	restores   int
}

func NewManager(store StatusStore) *Manager {
	return &Manager{
		store:      store,
		absorbers:  make(map[string]*absorber.Device),
		alarms:     make(map[string]string),
		history:    NewHistory(),
		escalation: NewEscalation(),
	}
}

func (m *Manager) Register(device *absorber.Device) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.absorbers[device.ID()] = device
}

func (m *Manager) device(deviceID string) *absorber.Device {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.absorbers[deviceID]
}

func (m *Manager) Engage(deviceID string) error {
	device := m.device(deviceID)
	if device == nil {
		return model.ErrDeviceFault
	}
	if err := device.Engage(); err != nil {
		m.raise(deviceID, "engage failure")
		return fmt.Errorf("engage %s: %w", deviceID, err)
	}
	m.clear(deviceID)
	return nil
}

func (m *Manager) RestoreDevice(deviceID string) error {
	device := m.device(deviceID)
	if device == nil {
		return model.ErrDeviceFault
	}
	if err := m.writebackWithRetry(deviceID, model.StatusRecovered); err != nil {
		return fmt.Errorf("restore writeback: %w", err)
	}
	if err := device.Restore(); err != nil {
		return err
	}
	m.restores++
	m.clear(deviceID)
	return nil
}

func (m *Manager) Recover(snapshot model.SystemSnapshot) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.revision = snapshot.Revision
	for _, snap := range snapshot.Devices {
		device, ok := m.absorbers[snap.ID]
		if !ok {
			continue
		}
		device.ApplySnapshot(snap)
		if snap.Status == model.StatusFault {
			m.alarms[snap.ID] = "device fault"
		} else {
			delete(m.alarms, snap.ID)
		}
	}
	return nil
}

func (m *Manager) AlarmCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.alarms)
}

func (m *Manager) Revision() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.revision
}

func (m *Manager) HasAlarm(deviceID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.alarms[deviceID]
	return ok
}

func (m *Manager) raise(deviceID, reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alarms[deviceID] = reason
	m.history.Record(deviceID, reason)
	m.escalation.Set(deviceID, LevelWarning)
}

func (m *Manager) clear(deviceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.alarms, deviceID)
	m.history.Clear(deviceID)
}

func (m *Manager) History() *History {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.history
}

func (m *Manager) Escalation() *Escalation {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.escalation
}

func (m *Manager) writebackWithRetry(deviceID string, status model.DeviceStatus) error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		if err := m.store.Writeback(deviceID, status); err == nil {
			return nil
		} else {
			last = err
		}
	}
	return last
}

func (m *Manager) RestoreCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.restores
}
