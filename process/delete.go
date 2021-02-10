package process

import (
	"fmt"
	"os"
)

type deleteaction struct{}

func deletefile(path string) error {
	return os.Remove(path)
}

func (a deleteaction) Process(path string, info os.FileInfo, dests []string) error {
	for _, d := range dests {
		err := deletefile(d)
		if err != nil {
			return fmt.Errorf("deleting file %s. %w", d, err)
		}
	}

	return nil
}

func (a deleteaction) RequireFileExist() bool {
	return false
}
