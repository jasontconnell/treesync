package conf

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/jasontconnell/conf"
)

var NoTreesyncErr error = errors.New("No Treesync config found")

type Config struct {
	Root          string
	RootFolders   []string `json:"rootFolders"`
	AlwaysExclude []string `json:"alwaysExclude"`
}

func GetExcludeMap(exclude string, excludes []string) map[string]bool {
	m := make(map[string]bool)
	m[exclude] = true
	for _, e := range excludes {
		m[e] = true
	}
	return m
}

func FindRoot(wd, filename string) (Config, error) {
	tmp := wd

	var configFile string
	found := false
	for !found {
		fullPath := filepath.Join(tmp, filename)
		f, err := os.Stat(fullPath)
		found = err == nil && !f.IsDir()
		if found {
			configFile = fullPath
			break
		}

		parent := filepath.Dir(tmp)
		if tmp == parent {
			break
		}
		tmp = parent
	}
	if !found {
		return Config{}, NoTreesyncErr
	}
	cfg := LoadConfig(configFile)
	cfg.Root = filepath.Dir(configFile)

	return cfg, nil
}

func LoadConfig(filename string) Config {
	cfg := Config{}
	conf.LoadConfig(filename, &cfg)
	return cfg
}
