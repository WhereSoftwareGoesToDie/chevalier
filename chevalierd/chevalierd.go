package main

import (
	"github.com/anchor/picolog"
	"code.google.com/p/gcfg"
	"os"
	"flag"
	"log"
)

const (
	Version = "0.2.1"
)

var Logger *picolog.Logger

func main() {
	var cfg Config
	configFile := flag.String("cfg", "/etc/chevalier.gcfg", "Path to configuration file. This file should be in gcfg[0] format. [0] https://code.google.com/p/gcfg/")
	indexerMode := flag.Bool("index", false, "Start indexer mode.")
	indexOnce := flag.Bool("index-once", false, "Run indexer once and then exit.")
	readerMode := flag.Bool("read", false, "Start reader mode. Reader mode will invoke the indexer once on startup; use -no-index to disable.")
	noIndex := flag.Bool("no-index", false, "Do not index once at startup when started in reader mode.")
	logFile := flag.String("log-file", "", "If set, log to this file rather than stdout.")
	flag.Parse()
	err := gcfg.ReadFileInto(&cfg, *configFile)
	if err != nil {
		log.Fatalf("Could not read config file at %v: %v", *configFile, err)
	}
	logLevel, err := picolog.ParseLogLevel(cfg.Chevalier.LogLevel)
	if err != nil {
		log.Fatalf("Could not parse log level: %v", err)
	}
	logStream := os.Stdout
	if *logFile != "" {
		logStream, err = os.OpenFile(*logFile, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0660);
		if err != nil {
			log.Fatalf("Could not open log file %v: %v", *logFile, err)
		}
	}
	Logger = picolog.NewLogger(logLevel, "chevalier", logStream)
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
