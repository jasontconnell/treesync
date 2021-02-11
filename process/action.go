package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Action interface {
	Process(path string, info os.FileInfo, dests []string) error
	RequireFileExist() bool
}

type DirAction interface {
	Action
	ProcessDir(path string, info os.FileInfo, dests []string) error
}

func getAction(name string) Action {
	var a Action
	switch name {
	case "copy":
		a = copyaction{}
	case "delete":
		a = deleteaction{}
	case "sync":
		a = syncaction{}
	default:
		panic("not defined")
	}

	return a
}

func getFiles(treesyncroot, wd, file string, roots, excludeMap map[string]bool) (string, []string) {
	fullPath := filepath.Clean(filepath.Join(wd, file))
	if filepath.IsAbs(file) {
		fullPath = file
	}

	curroot := strings.TrimPrefix(string(fullPath[len(treesyncroot):]), string(filepath.Separator))
	parts := strings.Split(curroot, string(filepath.Separator))
	curroot = parts[0]
	currootabs := filepath.Join(treesyncroot, curroot)
	relPath := strings.TrimPrefix(fullPath, currootabs)

	paths := []string{}

	for r := range roots {
		if r == curroot {
			continue
		}
		if _, ok := excludeMap[r]; ok {
			continue
		}
		p := filepath.Join(treesyncroot, r, relPath)
		paths = append(paths, p)
	}
	return fullPath, paths
}

func Process(action string, curdir, treesyncroot, file string, excludeMap map[string]bool, roots map[string]bool) error {
	a := getAction(action)
	cur, files := getFiles(treesyncroot, curdir, file, roots, excludeMap)

	stat, err := os.Stat(cur)
	if err != nil && a.RequireFileExist() {
		return fmt.Errorf("file or folder doesn't exist and %s requires it: %s %s %w", action, curdir, file, err)
	}

	isdir := stat != nil && stat.IsDir()

	da, candir := a.(DirAction)
	if candir && isdir {
		err = da.ProcessDir(cur, stat, files)
	} else if !isdir {
		err = a.Process(cur, stat, files)
	} else {
		err = fmt.Errorf("path is a directory and %s action can't process directories", action)
	}

	if err != nil {
		return fmt.Errorf("error processing action %s on file %s. %w", action, cur, err)
	}

	return err
}
