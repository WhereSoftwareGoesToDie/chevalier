package main

import (
	"code.google.com/p/gcfg"
	"github.com/anchor/picolog"
	"github.com/anchor/chevalier"

	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var Version = chevalier.Version

var Logger *picolog.Logger

// handleInterrupts will ignore HUP and WINCH and terminate on QUIT, INT
// and TERM. Does not do anything special at this stage, just run it in
// a goroutine at startup after the global logger has been initialized.
func handleInterrupts() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGWINCH)
	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGHUP, syscall.SIGWINCH:
			Logger.Infof("Ignoring %v.", sig)
		case syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
			Logger.Infof("Got %v, exiting.", sig)
			os.Exit(0)
		default:
			Logger.Fatalf("Got unexpected signal %v, panicking.")
		}
	}
}

func main() {
	var cfg Config
	configFile := flag.String("cfg", "/etc/chevalier.gcfg", "Path to configuration file. This file should be in gcfg[0] format. [0] https://code.google.com/p/gcfg/")
	indexerMode := flag.Bool("index", false, "Start indexer mode.")
	indexOnce := flag.Bool("index-once", false, "Run indexer once and then exit.")
	readerMode := flag.Bool("read", false, "Start reader mode. Reader mode will invoke the indexer once on startup; use -no-index to disable.")
	noIndex := flag.Bool("no-index", false, "Do not index once at startup when started in reader mode.")
	logFile := flag.String("log-file", "", "If set, log to this file rather than stdout.")
	debug := flag.Bool("debug", false, "Enable debug logging (overrides configured log level).")
	flag.Parse()
	err := gcfg.ReadFileInto(&cfg, *configFile)
	if err != nil {
		log.Fatalf("Could not read config file at %v: %v", *configFile, err)
	}
	logLevel, err := picolog.ParseLogLevel(cfg.Chevalier.LogLevel)
	if err != nil {
		log.Fatalf("Could not parse log level: %v", err)
	}
	if *debug {
		logLevel = picolog.LogDebug
	}
	logStream := os.Stdout
	if *logFile != "" {
		logStream, err = os.OpenFile(*logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			log.Fatalf("Could not open log file %v: %v", *logFile, err)
		}
	}
	Logger = picolog.NewLogger(logLevel, "chevalier", logStream)
	go handleInterrupts()
	if *indexerMode {
		RunIndexer(cfg)
	} else if *readerMode {
		if !*noIndex {
			go RunIndexerOnce(cfg)
		}
		go RunStatus(cfg)
		RunReader(cfg)
	} else if *indexOnce {
		RunIndexerOnce(cfg)
	} else {
		Logger.Fatalf("Must specify one of -index, -read or -index-once.")
	}
}
