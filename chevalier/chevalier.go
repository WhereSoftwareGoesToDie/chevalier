package main

import (
	zmq "github.com/pebbe/zmq4"
	"flag"
	"log"
	"code.google.com/p/gcfg"
)

func main() {
	var cfg Config
	configFile := flag.String("cfg", "/etc/chevalier.gcfg", "Path to configuration file. This file should be in gcfg[0] format. [0] https://code.google.com/p/gcfg/")
	flag.Parse()
	err := gcfg.ReadFileInto(&cfg, *configFile)
	if err != nil {
		log.Fatalf("Could not read config file at %v: %v", *configFile, err)
	}
	sock, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		log.Fatalf("Could not initialize listen socket: %v", err)
	}
	err = sock.Bind(cfg.Chevalier.ListenAddress)
	if err != nil {
		log.Fatalf("Could not listen on %v: %v", cfg.Chevalier.ListenAddress, err)
	}
	for true {
		_, err = sock.RecvMessageBytes(0)
	}
}
