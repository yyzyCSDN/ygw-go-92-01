package record

import (
	"fmt"
	"os"
	"path/filepath"
)

type Verifier struct {
	dir string
}

func NewVerifier(dir string) *Verifier {
	return &Verifier{dir: dir}
}

func (v *Verifier) Verify() (int, error) {
	files, err := filepath.Glob(filepath.Join(v.dir, "*.log"))
	if err != nil {
		return 0, err
	}
	total := 0
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return 0, fmt.Errorf("read %s: %w", path, err)
		}
		total += len(data)
	}
	return total, nil
}
