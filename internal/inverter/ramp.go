package inverter

type Ramp struct {
	rate    float64
	current float64
	target  float64
}

func NewRamp(rate float64) *Ramp {
	return &Ramp{rate: rate}
}

func (r *Ramp) SetTarget(target float64) {
	r.target = target
}

func (r *Ramp) Step() float64 {
	delta := r.target - r.current
	if delta > r.rate {
		delta = r.rate
	} else if delta < -r.rate {
		delta = -r.rate
	}
	r.current += delta
	return r.current
}

func (r *Ramp) Current() float64 {
	return r.current
}

func (r *Ramp) Settled() bool {
	delta := r.target - r.current
	return delta < r.rate && delta > -r.rate
}
