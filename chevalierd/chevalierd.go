package main

import (
	"github.com/anchor/picolog"
	"code.google.com/p/gcfg"
	"os"
	"flag"
	"log"
)

const (
	Version = "0.0.1"
)

var Logger *picolog.Logger

func main() {
	var cfg Config
	configFile := flag.String("cfg", "/etc/chevalier.gcfg", "Path to configuration file. This file should be in gcfg[0] format. [0] https://code.google.com/p/gcfg/")
	indexerMode := flag.Bool("index", false, "Start indexer mode.")
	readerMode := flag.Bool("read", false, "Start reader mode.")
	flag.Parse()
	err := gcfg.ReadFileInto(&cfg, *configFile)
	if err != nil {
		log.Fatalf("Could not read config file at %v: %v", *configFile, err)
	}
	logLevel, err := picolog.ParseLogLevel(cfg.Chevalier.LogLevel)
	if err != nil {
		log.Fatalf("Could not parse log level: %v", err)
	}
	Logger = picolog.NewLogger(logLevel, "chevalier", os.Stdout)
	if *indexerMode {
		runIndexer(cfg)
	} else if *readerMode {
		runReader(cfg)
	} else {
		Logger.Fatalf("Must specify either -index or -read.")
	}
}
