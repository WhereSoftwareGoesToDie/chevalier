package main

import (
	"github.com/anchor/picolog"
	"code.google.com/p/gcfg"
	"os"
	"flag"
	"log"
)

const (
	Version = "0.2.0"
)

var Logger *picolog.Logger

func main() {
	var cfg Config
	configFile := flag.String("cfg", "/etc/chevalier.gcfg", "Path to configuration file. This file should be in gcfg[0] format. [0] https://code.google.com/p/gcfg/")
	indexerMode := flag.Bool("index", false, "Start indexer mode.")
	indexOnce := flag.Bool("index-once", false, "Run indexer once and then exit.")
	readerMode := flag.Bool("read", false, "Start reader mode. Reader mode will invoke the indexer once on startup; use -no-index to disable.")
	noIndex := flag.Bool("no-index", false, "Do not index once at startup when started in reader mode.")
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
		RunIndexer(cfg)
	} else if *readerMode {
		if !*noIndex {
			go RunIndexerOnce(cfg)
		}
		RunReader(cfg)
	} else if *indexOnce {
		RunIndexerOnce(cfg)
	} else {
		Logger.Fatalf("Must specify one of -index, -read or -index-once.")
	}
}
