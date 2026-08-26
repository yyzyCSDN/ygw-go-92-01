package coop

import "regenbrake/internal/model"

func (c *Coordinator) Rebalance(limit float64) []model.Allocation {
	c.mu.Lock()
	defer c.mu.Unlock()
	current := c.total
	if current <= limit {
		return nil
	}
	scale := limit / current
	out := make([]model.Allocation, 0, len(c.alloc))
	c.total = 0
	for id, item := range c.alloc {
		item.Capacity *= scale
		c.alloc[id] = item
		c.total += item.Capacity
		out = append(out, item)
	}
	return out
}

func (c *Coordinator) CapacityPerTrain() map[string]float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]float64, len(c.alloc))
	for id, item := range c.alloc {
		out[id] = item.Capacity
	}
	return out
}
