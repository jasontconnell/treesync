package process

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type copyaction struct{}

func copyfile(path, dest string) error {
	dir, name := filepath.Split(dest)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating directory structure for %s. %w", dest, err)
	}

	src, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening file for read to copy: %s. %w", name, err)
	}
	defer src.Close()

	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating file to write: %s. %w", dest, err)
	}
	defer f.Close()

	_, err = io.Copy(f, src)
	return err
}

func (a copyaction) Process(path string, info os.FileInfo, dests []string) error {
	for _, d := range dests {
		err := copyfile(path, d)
		if err != nil {
			return fmt.Errorf("copy file failed: %s -> %s. %w", path, d, err)
		}
	}

	return nil
}

func (a copyaction) ProcessDir(path string, info os.FileInfo, dests []string) error {
	copyRelPaths := []string{}
	filepath.Walk(path, func(cpath string, inf os.FileInfo, err error) error {
		if inf.IsDir() {
			return nil
		}

		rel := strings.TrimPrefix(cpath, path)
		copyRelPaths = append(copyRelPaths, rel)

		return nil
	})

	for _, d := range dests {
		for _, c := range copyRelPaths {
			sfull := filepath.Join(path, c)
			dfull := filepath.Join(d, c)
			err := copyfile(sfull, dfull)
			if err != nil {
				return fmt.Errorf("copy file failed. %s -> %s. %w", sfull, dfull, err)
			}
		}
	}

	return nil
}

func (a copyaction) RequireFileExist() bool {
	return true
}
