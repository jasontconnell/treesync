package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jasontconnell/treesync/conf"
	"github.com/jasontconnell/treesync/process"
)

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

	wd, err := os.Getwd()
	cfg, err := conf.FindRoot(wd, *cfgfile)
	if err == conf.NoTreesyncErr {
		log.Println("no tree sync root config file found", *cfgfile)
	}

	emap := conf.GetStringMap(strings.Split(*exclude, ","), cfg.AlwaysExclude)
	var incmap map[string]bool
	if *include != "" {
		incmap = conf.GetStringMap(strings.Split(*include, ","))
	} else {
		incmap = conf.GetStringMap(cfg.RootFolders)
	}

	err = process.Process(*action, wd, cfg.Root, file, emap, incmap)
	if err != nil {
		log.Fatal("error processing", err)
	}
	log.Println("Success. Time:", time.Since(start))
}
