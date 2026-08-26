package absorber

type Metrics struct {
	Engaged    int
	Discharged int
	Recovered  int
	Faults     int
}

func (m *Metrics) RecordEngage()    { m.Engaged++ }
func (m *Metrics) RecordDischarge() { m.Discharged++ }
func (m *Metrics) RecordRecover()   { m.Recovered++ }
func (m *Metrics) RecordFault()     { m.Faults++ }
