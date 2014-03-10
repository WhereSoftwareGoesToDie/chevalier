package main

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/anchor/chevalier"
	"github.com/anchor/picolog"
)

var IndexerLogger *picolog.Logger

func fullUpdate(w *chevalier.ElasticsearchWriter, endpoint string, origins []string) {
	indexed := 0
	for _, o := range origins {
		burst, err := chevalier.GetContents(endpoint, o)
		if err != nil {
			IndexerLogger.Errorf("Could not read contents for origin %v: %v", o, err)
			continue
		}
		for _, s := range burst.Sources {
			err = w.Write(s)
			if err != nil {
				IndexerLogger.Errorf("Could not index source: %v", err)
			} else {
				indexed += 1
			}
		}
	}
	IndexerLogger.Infof("Indexed %v sources.", indexed)
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
	IndexerLogger = Logger.NewSubLogger("indexer")
	writer := chevalier.NewElasticsearchWriter(cfg.Elasticsearch.Host, cfg.Elasticsearch.MaxConns, cfg.Elasticsearch.RetrySeconds, cfg.Elasticsearch.DataType, cfg.Elasticsearch.DataType)
	fullUpdate(writer, cfg.Vaultaire.ReadEndpoint, cfg.Vaultaire.Origins)
}
