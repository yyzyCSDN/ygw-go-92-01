package record

import "time"

type Stats struct {
	Records   int
	Bytes     int
	Rotations int
	StartedAt time.Time
	LastWrite time.Time
}

type Counter struct {
	records   int
	bytes     int
	rotations int
	startedAt time.Time
	lastWrite time.Time
}

func NewCounter() *Counter {
	now := time.Now().UTC()
	return &Counter{startedAt: now, lastWrite: now}
}

func (c *Counter) RecordWrite(data []byte) {
	c.records++
	c.bytes += len(data)
	c.lastWrite = time.Now().UTC()
}

func (c *Counter) RecordRotate() {
	c.rotations++
}

func (c *Counter) Snapshot() Stats {
	return Stats{
		Records:   c.records,
		Bytes:     c.bytes,
		Rotations: c.rotations,
		StartedAt: c.startedAt,
		LastWrite: c.lastWrite,
	}
}
