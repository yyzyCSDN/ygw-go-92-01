package absorber

import "sort"

type Ledger struct {
	entries map[string]float64
}

func NewLedger() *Ledger {
	return &Ledger{entries: make(map[string]float64)}
}

func (l *Ledger) Set(trainID string, capacity float64) {
	l.entries[trainID] = capacity
}

func (l *Ledger) Get(trainID string) float64 {
	return l.entries[trainID]
}

func (l *Ledger) Remove(trainID string) {
	delete(l.entries, trainID)
}

func (l *Ledger) Total() float64 {
	total := 0.0
	for _, capacity := range l.entries {
		total += capacity
	}
	return total
}

func (l *Ledger) Keys() []string {
	out := make([]string, 0, len(l.entries))
	for trainID := range l.entries {
		out = append(out, trainID)
	}
	sort.Strings(out)
	return out
}
