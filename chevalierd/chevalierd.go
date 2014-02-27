package main

import (
	"github.com/anchor/picolog"
	"log/syslog"
	"code.google.com/p/gcfg"
	"os"
	"flag"
)

var Logger *picolog.Logger

func main() {
	var cfg Config
	Logger = picolog.NewLogger(syslog.LOG_DEBUG, "chevalier", os.Stdout)
	configFile := flag.String("cfg", "/etc/chevalier.gcfg", "Path to configuration file. This file should be in gcfg[0] format. [0] https://code.google.com/p/gcfg/")
	indexerMode := flag.Bool("index", false, "Start indexer mode.")
	readerMode := flag.Bool("read", false, "Start reader mode.")
	flag.Parse()
	err := gcfg.ReadFileInto(&cfg, *configFile)
	if err != nil {
		Logger.Fatalf("Could not read config file at %v: %v", *configFile, err)
	}
	if *indexerMode {
		runIndexer(cfg)
	} else if *readerMode {
		runReader(cfg)
	} else {
		Logger.Fatalf("Must specify either -index or -read.")
	}
}
