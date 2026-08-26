package record

import "time"

type Entry struct {
	At      time.Time
	Message string
	Digest  string
}

type History struct {
	entries []Entry
}

func NewHistory() *History {
	return &History{entries: make([]Entry, 0)}
}

func (h *History) Add(entry Entry) {
	h.entries = append(h.entries, entry)
}

func (h *History) Len() int {
	return len(h.entries)
}

func (h *History) Last() (Entry, bool) {
	if len(h.entries) == 0 {
		return Entry{}, false
	}
	return h.entries[len(h.entries)-1], true
}
