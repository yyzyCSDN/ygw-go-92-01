package alarm

import (
	"sort"
	"time"
)

type Event struct {
	At       time.Time
	DeviceID string
	Reason   string
	Cleared  bool
}

type History struct {
	events []Event
}

func NewHistory() *History {
	return &History{events: make([]Event, 0)}
}

func (h *History) Record(deviceID, reason string) {
	h.events = append(h.events, Event{At: time.Now().UTC(), DeviceID: deviceID, Reason: reason})
}

func (h *History) Clear(deviceID string) {
	for index := len(h.events) - 1; index >= 0; index-- {
		if h.events[index].DeviceID == deviceID && !h.events[index].Cleared {
			h.events[index].Cleared = true
			return
		}
	}
}

func (h *History) OpenEvents() []Event {
	out := make([]Event, 0)
	for _, event := range h.events {
		if !event.Cleared {
			out = append(out, event)
		}
	}
	return out
}

func (h *History) ByDevice(deviceID string) []Event {
	out := make([]Event, 0)
	for _, event := range h.events {
		if event.DeviceID == deviceID {
			out = append(out, event)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}
