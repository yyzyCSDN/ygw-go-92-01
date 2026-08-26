package record

import (
	"fmt"
	"path/filepath"
	"sync"
)

type Journal struct {
	mu        sync.Mutex
	dir       string
	writer    *FileWriter
	sequence  int
	counter   *Counter
	history   *History
	closed    bool
	rotations int
}

func NewJournal(dir string) *Journal {
	return &Journal{dir: dir, writer: &FileWriter{}, counter: NewCounter(), history: NewHistory()}
}

func (j *Journal) Append(entry string) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.writer.file == nil {
		if err := j.writer.Open(j.nextPath()); err != nil {
			return err
		}
	}
	data := []byte(entry + "\n")
	j.counter.RecordWrite(data)
	j.history.Add(Entry{Message: entry, Digest: Digest(data)})
	return j.writer.Write(data)
}

func (j *Journal) Rotate() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.closed {
		return fmt.Errorf("journal already closed")
	}
	j.sequence++
	j.counter.RecordRotate()
	j.rotations++
	if j.writer.file != nil {
		if err := j.writer.Close(); err != nil {
			return err
		}
	}
	return j.writer.Open(j.nextPath())
}

func (j *Journal) Close() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.closed {
		return nil
	}
	j.closed = true
	return j.writer.Close()
}

func (j *Journal) Sequence() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.sequence
}

func (j *Journal) Stats() Stats {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.counter.Snapshot()
}

func (j *Journal) LastEntry() (Entry, bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.history.Last()
}

func (j *Journal) EntryCount() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.history.Len()
}

func (j *Journal) Closed() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.closed
}

func (j *Journal) RotationCount() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.rotations
}

func (j *Journal) nextPath() string {
	return filepath.Join(j.dir, fmt.Sprintf("run-%06d.log", j.sequence))
}
