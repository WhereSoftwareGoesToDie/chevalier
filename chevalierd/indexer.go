package main

import (
	zmq "github.com/pebbe/zmq4"
	_ "github.com/anchor/chevalier"
)

func fullUpdate(endpoint string) error {
	sock, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		return err
	}
	err = sock.Connect(endpoint)
	if err != nil {
		return err
	}
	_, err = sock.Send("", 0)
	if err != nil {
		return err
	}
	return nil
}

func subscribeUpdate(endpoint string) error {
	sock, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		return err
	}
	sock.SetSubscribe("")
	sock.Connect(endpoint)
	return nil
}

func runIndexer(cfg Config) {
	Logger.Infof("Starting chevalierd %v in indexer mode.", Version)
	go fullUpdate(cfg.Vaultaire.ReadEndpoint)
}
