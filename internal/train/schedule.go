package train

import "sort"

type Schedule struct {
	arrivals map[string]int
}

func NewSchedule() *Schedule {
	return &Schedule{arrivals: make(map[string]int)}
}

func (s *Schedule) Add(trainID string, station int) {
	s.arrivals[trainID] = station
}

func (s *Schedule) Station(trainID string) int {
	return s.arrivals[trainID]
}

func (s *Schedule) TrainsAt(station int) []string {
	out := make([]string, 0)
	for trainID, current := range s.arrivals {
		if current == station {
			out = append(out, trainID)
		}
	}
	sort.Strings(out)
	return out
}
