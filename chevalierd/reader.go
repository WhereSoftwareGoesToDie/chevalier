package main

import (
	"github.com/anchor/chevalier"
	zmq "github.com/pebbe/zmq4"
)

func handleRequest(sock *zmq.Socket, engine *chevalier.QueryEngine) error {
	msg, err := sock.RecvBytes(0)
	Logger.Debugf("Got a request!")
	req, err := chevalier.UnmarshalSourceRequest(msg)
	if err != nil {
		Logger.Warningf("Failed to unmarshal request: %v", err)
		return err
	}
	Logger.Debugf("%v", req)
	results, err := engine.GetSources(req)
	if err != nil {
		Logger.Errorf("Error querying Elasticsearch: %v", err)
		return err
	}
	Logger.Debugf("Got result: %v", results)
	reply, err := chevalier.MarshalSourceBurst(results)
	if err != nil {
		Logger.Errorf("Error marshalling reply: %v", err)
	}
	_, err = sock.SendBytes(reply, 0)
	if err != nil {
		Logger.Errorf("Error sending response: %v", err)
	}
	return nil
}

func runReader(cfg Config) {
	Logger.Infof("Starting chevalierd %v in reader mode.", Version)
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
	for {
		err = reactor.Run(-1)
		if err != nil {
			Logger.Errorf("Restarting reactor: %v", err)
		}
	}
}
