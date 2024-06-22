package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jasontconnell/treesync/conf"
	"github.com/jasontconnell/treesync/process"
)

func setLog(wd, name string) {
	if name != "stdout" {
		// exe, _ := os.Executable()
		logfile := filepath.Join(wd, name)
		f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Println("couldn't open log file, using stdout")
			return
		}
		log.SetOutput(f)
	}
}

func main() {
	start := time.Now()
	cfgfile := flag.String("c", "treesync.json", "config file")
	action := flag.String("a", "", "action to perform: sync, copy, delete (sync works on folders, will replace tree)")
	exclude := flag.String("exc", "", "root folders to exclude from this action")
	include := flag.String("inc", "", "root folders to include in this action")
	flag.Parse()

	file := os.Args[len(os.Args)-1]

	if file == "" || *action == "" {
		flag.PrintDefaults()
		return
	}

	actions := map[string]bool{
		"sync":   true,
		"copy":   true,
		"delete": true,
	}

	if _, ok := actions[*action]; !ok {
		log.Fatalf("invalid action %s", *action)
	}

	usefile := filepath.IsAbs(file)
	wd, err := os.Getwd()
	if usefile {
		wd = filepath.Dir(file)
	}

	cfg, err := conf.FindRoot(wd, *cfgfile)

	if err != nil {
		setLog(cfg.Root, cfg.Log)
	}

	if err == conf.NoTreesyncErr {
		log.Fatal("no tree sync root config file found", wd, *cfgfile, file)
	}

	var excfinal []string
	exlist := strings.Split(*exclude, ",")
	for _, ex := range exlist {
		if list, ok := cfg.FolderGroups[ex]; ok {
			excfinal = append(excfinal, list...)
		} else {
			excfinal = append(excfinal, ex)
		}
	}
	emap := conf.GetStringMap(excfinal, cfg.AlwaysExclude)

	var incmap map[string]bool
	if *include != "" {
		incfinal := []string{}
		inclist := strings.Split(*include, ",")
		for _, inc := range inclist {
			if list, ok := cfg.FolderGroups[inc]; ok {
				incfinal = append(incfinal, list...)
			} else {
				incfinal = append(incfinal, inc)
			}
		}
		incmap = conf.GetStringMap(incfinal)
	} else {
		incmap = conf.GetStringMap(cfg.RootFolders)
	}

	err = process.Process(*action, wd, cfg.Root, file, emap, incmap)
	if err != nil {
		log.Fatal(os.Args, "\nerror processing - ", cfg.Root, "\n", file, "\n", err)
	}
	log.Println("Success. Time:", time.Since(start))
}
