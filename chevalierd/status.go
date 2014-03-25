package main

import (
	"github.com/anchor/chevalier"
	"github.com/anchor/picolog"
	"github.com/anchor/zmqutil"
	zmq "github.com/pebbe/zmq4"
)

var StatusLogger *picolog.Logger

func handleStatusRequest(sock *zmq.Socket, engine *chevalier.QueryEngine, origins []string) error {
	_, err := zmqutil.RetryRecvMessage(sock, 0)
	if err != nil {
		StatusLogger.Errorf("Error receiving status request: %v", err)
	}
	StatusLogger.Debugf("Got a status request!")
	results := engine.GetStatus(origins)
	StatusLogger.Debugf("Got result: %v", results)
	reply, err := results.Marshal()
	if err != nil {
		StatusLogger.Errorf("Error marshalling reply: %v", err)
	}
	_, err = sock.SendBytes(reply, 0)
	if err != nil {
		StatusLogger.Errorf("Error sending response: %v", err)
	}
	return nil
}

func RunStatus(cfg Config) {
	StatusLogger = Logger.NewSubLogger("status")
	sock, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		StatusLogger.Fatalf("Could not initialize listen socket: %v", err)
	}
	err = sock.Bind(cfg.Chevalier.StatusAddress)
	if err != nil {
		StatusLogger.Fatalf("Could not listen on %v: %v", cfg.Chevalier.ListenAddress, err)
	}
	engine := chevalier.NewQueryEngine(cfg.Elasticsearch.Host, cfg.Elasticsearch.Index, cfg.Elasticsearch.DataType, cfg.Elasticsearch.MetadataIndex)
	reactor := zmq.NewReactor()
	reactor.AddSocket(sock, zmq.POLLIN, func(e zmq.State) error { return handleStatusRequest(sock, engine, cfg.Vaultaire.Origins) })
	for {
		err = reactor.Run(-1)
		if err != nil {
			StatusLogger.Errorf("Restarting reactor: %v", err)
		}
	}
}
