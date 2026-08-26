package alarm

import "regenbrake/internal/model"

type MemoryStore struct {
	statuses map[string]model.DeviceStatus
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{statuses: make(map[string]model.DeviceStatus)}
}

func (s *MemoryStore) Writeback(id string, status model.DeviceStatus) error {
	s.statuses[id] = status
	return nil
}

func (s *MemoryStore) Status(id string) model.DeviceStatus {
	return s.statuses[id]
}
