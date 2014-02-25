package main

import (
	"code.google.com/p/gcfg"
	"flag"
	zmq "github.com/pebbe/zmq4"
	"github.com/anchor/picolog"
	"github.com/anchor/chevalier"
	"os"
	"log/syslog"
)

var Logger *picolog.Logger

func handleRequest(sock *zmq.Socket, engine *chevalier.QueryEngine) error {
	msg, err := sock.RecvBytes(0)
	Logger.Debugf("Got a request!")
	req, err := chevalier.UnmarshalSourceRequest(msg)
	if err != nil {
		Logger.Warningf("Failed to unmarshal request: %v", err)
	}
	Logger.Debugf("%v", req)
	results, err := engine.RunSourceRequest(req)
	if err != nil {
		Logger.Errorf("Error querying Elasticsearch: %v", err)
		return nil
	}
	sources := chevalier.FmtResult(results)
	Logger.Debugf("Got result: %v", sources)
	return nil
}

func main() {
	var cfg Config
	Logger = picolog.NewLogger(syslog.LOG_DEBUG, "chevalier", os.Stdout)
	configFile := flag.String("cfg", "/etc/chevalier.gcfg", "Path to configuration file. This file should be in gcfg[0] format. [0] https://code.google.com/p/gcfg/")
	flag.Parse()
	err := gcfg.ReadFileInto(&cfg, *configFile)
	if err != nil {
		Logger.Fatalf("Could not read config file at %v: %v", *configFile, err)
	}
	sock, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		Logger.Fatalf("Could not initialize listen socket: %v", err)
	}
	err = sock.Bind(cfg.Chevalier.ListenAddress)
	if err != nil {
		Logger.Fatalf("Could not listen on %v: %v", cfg.Chevalier.ListenAddress, err)
	}
	engine := chevalier.NewQueryEngine(cfg.Elasticsearch.Host, cfg.Elasticsearch.Index, cfg.Elasticsearch.DataType)
	reactor := zmq.NewReactor()
	reactor.AddSocket(sock, zmq.POLLIN, func(e zmq.State) error { return handleRequest(sock, engine) })
	err = reactor.Run(-1)
	if err != nil {
		Logger.Errorf("%v", err)
	}
}
