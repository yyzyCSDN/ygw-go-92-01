package voltage

type Hysteresis struct {
	upper float64
	lower float64
	high  bool
}

func NewHysteresis(upper, lower float64) *Hysteresis {
	return &Hysteresis{upper: upper, lower: lower}
}

func (h *Hysteresis) Update(voltage float64) bool {
	if voltage >= h.upper {
		h.high = true
	} else if voltage <= h.lower {
		h.high = false
	}
	return h.high
}

func (h *Hysteresis) High() bool {
	return h.high
}

func (h *Hysteresis) Band() (float64, float64) {
	return h.lower, h.upper
}
