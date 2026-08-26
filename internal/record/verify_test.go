package record_test

import (
	"os"
	"path/filepath"
	"testing"

	"regenbrake/internal/record"
)

func TestRunHandleClosed(t *testing.T) {
	dir := t.TempDir()
	journal := record.NewJournal(dir)
	for i := 0; i < 40; i++ {
		if err := journal.Append("sample"); err != nil {
			t.Fatal(err)
		}
		if err := journal.Rotate(); err != nil {
			t.Fatal(err)
		}
	}
	if err := journal.Close(); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		src := filepath.Join(dir, entry.Name())
		if err := os.Rename(src, src+".renamed"); err != nil {
			t.Fatalf("file handle still open for %s: %v", entry.Name(), err)
		}
	}
}
