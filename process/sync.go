package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type syncaction struct{}

var ErrSyncOnlyDirs = fmt.Errorf("sync action is only available on directories")

func (a syncaction) Process(path string, info os.FileInfo, dests []string) error {
	return ErrSyncOnlyDirs
}
func (a syncaction) RequireFileExist() bool {
	return true
}

func (a syncaction) ProcessDir(path string, info os.FileInfo, dests []string) error {
	syncpaths := []string{}

	filepath.Walk(path, func(spath string, inf os.FileInfo, err error) error {
		if inf.IsDir() {
			return nil
		}

		rel := strings.TrimPrefix(spath, path)
		syncpaths = append(syncpaths, rel)

		return nil
	})

	for _, d := range dests {
		err := os.RemoveAll(d)
		if err != nil {
			return fmt.Errorf("couldn't remove directory for sync: %s. %w", d, err)
		}
	}

	for _, d := range dests {
		for _, spath := range syncpaths {
			srcfile, destfile := filepath.Join(path, spath), filepath.Join(d, spath)
			err := copyfile(srcfile, destfile)
			if err != nil {
				return fmt.Errorf("couldn't sync file: %s -> %s. %w", srcfile, destfile, err)
			}
		}
	}

	return nil
}
