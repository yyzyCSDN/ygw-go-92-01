package absorber

import (
	"sort"
	"sync"
)

type Ledger struct {
	mu      sync.Mutex
	entries map[string]float64
}

func NewLedger() *Ledger {
	return &Ledger{entries: make(map[string]float64)}
}

func (l *Ledger) Set(trainID string, capacity float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries[trainID] = capacity
}

func (l *Ledger) Get(trainID string) float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.entries[trainID]
}

func (l *Ledger) Remove(trainID string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, trainID)
}

func (l *Ledger) Total() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	total := 0.0
	for _, capacity := range l.entries {
		total += capacity
	}
	return total
}

func (l *Ledger) Keys() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]string, 0, len(l.entries))
	for trainID := range l.entries {
		out = append(out, trainID)
	}
	sort.Strings(out)
	return out
}
