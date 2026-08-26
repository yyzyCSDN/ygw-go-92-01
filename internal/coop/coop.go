package coop

import (
	"sort"
	"sync"

	"regenbrake/internal/absorber"
	"regenbrake/internal/model"
	"regenbrake/internal/train"
)

type Coordinator struct {
	mu       sync.Mutex
	reg      *train.Registry
	abs      *absorber.Device
	state    model.CoopState
	alloc    map[string]model.Allocation
	total    float64
	released map[string]bool
	misses   int
	stale    int
	missing  map[string]bool
}

func NewCoordinator(reg *train.Registry, abs *absorber.Device) *Coordinator {
	return &Coordinator{
		reg:      reg,
		abs:      abs,
		state:    model.CoopPlanning,
		alloc:    make(map[string]model.Allocation),
		released: make(map[string]bool),
		missing:  make(map[string]bool),
	}
}

func (c *Coordinator) Allocate(reports []model.BrakingReport) ([]model.Allocation, error) {
	c.state = model.CoopAllocating
	out := make([]model.Allocation, 0, len(reports))
	for _, report := range reports {
		current, ok := c.reg.Lookup(report.TrainID)
		if !ok {
			c.misses++
			c.missing[report.TrainID] = true
			continue
		}
		if !current.Braking {
			continue
		}
		item := model.Allocation{TrainID: report.TrainID, Capacity: report.Capacity, Power: report.Power}
		c.alloc[report.TrainID] = item
		c.total += report.Capacity
		c.abs.AddCapacity(report.Capacity)
		out = append(out, item)
	}
	c.state = model.CoopSettled
	return out, nil
}

func (c *Coordinator) Allocation(trainID string) (model.Allocation, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.alloc[trainID]
	return item, ok
}

func (c *Coordinator) Total() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.total
}

func (c *Coordinator) State() model.CoopState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

func (c *Coordinator) OnTrainBrakingChange(trainID string, braking bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if braking {
		return
	}
	c.releaseTrainLocked(trainID)
}

func (c *Coordinator) releaseTrainLocked(trainID string) {
	item, ok := c.alloc[trainID]
	if !ok {
		return
	}
	if c.released[trainID] {
		c.stale++
		return
	}
	c.total -= item.Capacity
	c.abs.AddCapacity(-item.Capacity)
	c.released[trainID] = true
	delete(c.alloc, trainID)
}

func (c *Coordinator) Allocations() []model.Allocation {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]model.Allocation, 0, len(c.alloc))
	for _, item := range c.alloc {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].TrainID < out[j].TrainID })
	return out
}

func (c *Coordinator) MissCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.misses
}

func (c *Coordinator) StaleReleaseCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stale
}

func (c *Coordinator) MissingTrains() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]string, 0, len(c.missing))
	for trainID := range c.missing {
		out = append(out, trainID)
	}
	sort.Strings(out)
	return out
}
