package main

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/anchor/chevalier"
	"github.com/anchor/picolog"
	"syscall"
)

var IndexerLogger *picolog.Logger

func originUpdate(w *chevalier.ElasticsearchWriter, endpoint string, origin string) {
	indexed := 0
	IndexerLogger.Infof("Requesting sources for origin %v.", origin)
	// We want to retry if we get interrupted.
	var err error
	var burst *chevalier.DataSourceBurst
	err = syscall.EAGAIN
	for err == syscall.EAGAIN || err == syscall.EINTR {
		burst, err = chevalier.GetContents(endpoint, origin)
	}
	if err != nil {
		IndexerLogger.Errorf("Could not read contents for origin %v: %v", origin, err)
		return
	}
	for _, s := range burst.Sources {
		err = w.Write(origin, s)
		if err != nil {
			IndexerLogger.Errorf("Could not index source: %v", err)
		} else {
			indexed += 1
		}
	}
	IndexerLogger.Infof("Indexed %v sources for origin %v.", indexed, origin)
}

func fullUpdate(w *chevalier.ElasticsearchWriter, endpoint string, origins []string) {
	for _, o := range origins {
		go originUpdate(w, endpoint, o)
	}
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

func getElasticsearchWriter(cfg Config) *chevalier.ElasticsearchWriter {
	writer := chevalier.NewElasticsearchWriter(cfg.Elasticsearch.Host, cfg.Elasticsearch.MaxConns, cfg.Elasticsearch.RetrySeconds, cfg.Elasticsearch.Index, cfg.Elasticsearch.DataType)
	return writer
}

func RunIndexerOnce(cfg Config) {
	IndexerLogger = Logger.NewSubLogger("indexer")
	IndexerLogger.Infof("Starting single indexer run.")
	writer := getElasticsearchWriter(cfg)
	fullUpdate(writer, cfg.Vaultaire.ReadEndpoint, cfg.Vaultaire.Origins)
}

func RunIndexer(cfg Config) {
	Logger.Infof("Starting chevalierd %v in indexer mode.", Version)
	IndexerLogger = Logger.NewSubLogger("indexer")
	writer := getElasticsearchWriter(cfg)
	for {
		IndexerLogger.Infof("Starting run.")
		fullUpdate(writer, cfg.Vaultaire.ReadEndpoint, cfg.Vaultaire.Origins)
	}
}
