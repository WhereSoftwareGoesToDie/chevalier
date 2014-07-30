package main

import (
	"github.com/anchor/chevalier"
	"github.com/anchor/picolog"
	zmq "github.com/pebbe/zmq4"
	"syscall"
	"time"
)

var IndexerLogger *picolog.Logger

func originUpdate(sem chan bool, res chan uint64, w *chevalier.ElasticsearchWriter, endpoint, origin string) {
	// Block until there's an available slot.
	<-sem
	indexed := uint64(0)
	IndexerLogger.Infof("Requesting sources for origin %v.", origin)
	// We want to retry if we get interrupted.
	var err error
	var sources []*chevalier.ElasticsearchSource
	err = syscall.EAGAIN
	for err == syscall.EAGAIN || err == syscall.EINTR {
		sources, err = chevalier.GetContents(endpoint, origin)
	}
	if err != nil {
		IndexerLogger.Errorf("Could not read contents for origin %v: %v", origin, err)
		return
	}
	for _, s := range sources {
		err = w.Write(origin, s)
		IndexerLogger.Debugf("Writing source %v for origin %s.", s, origin)
		if err != nil {
			IndexerLogger.Errorf("Could not index source: %v", err)
		} else {
			indexed += 1
		}
	}
	IndexerLogger.Infof("Indexed %v sources for origin %v.", indexed, origin)
	err = w.UpdateOrigin(origin, indexed)
	if err != nil {
		IndexerLogger.Errorf("Could not update origin %v: %v", origin, err)
	} else {
		IndexerLogger.Infof("Updated metdata for origin %v.", origin)
	}
	// Seed semaphore for the next goroutine.
	sem <- true
	res <- indexed
}

func fullUpdate(w *chevalier.ElasticsearchWriter, endpoint string, origins []string, parallelism uint) {
	total := uint64(0)
	output := make(chan uint64, 0)
	semaphore := make(chan bool, parallelism)
	for i := uint(0); i < parallelism; i++ {
		semaphore <- true
	}
	for _, o := range origins {
		go originUpdate(semaphore, output, w, endpoint, o)
	}
	for _, _ = range origins {
		indexed := <-output
		total += indexed
	}
	IndexerLogger.Infof("Indexer run finished; indexed %v sources in total.", total)
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
	writer := chevalier.NewElasticsearchWriter(cfg.Elasticsearch.Host, cfg.Elasticsearch.MaxConns, cfg.Elasticsearch.RetrySeconds, cfg.Elasticsearch.Index, cfg.Elasticsearch.MetadataIndex, cfg.Elasticsearch.DataType)
	return writer
}

func RunIndexerOnce(cfg Config) {
	IndexerLogger = Logger.NewSubLogger("indexer")
	IndexerLogger.Infof("Starting single indexer run.")
	writer := getElasticsearchWriter(cfg)
	fullUpdate(writer, cfg.Vaultaire.ContentsEndpoint, cfg.Vaultaire.Origins, cfg.Indexer.Parallelism)
}

func RunIndexer(cfg Config) {
	Logger.Infof("Starting chevalierd %v in indexer mode. Index interval: %v seconds.", Version, cfg.Indexer.IndexInterval)
	IndexerLogger = Logger.NewSubLogger("indexer")
	writer := getElasticsearchWriter(cfg)
	for {
		IndexerLogger.Infof("Starting run.")
		fullUpdate(writer, cfg.Vaultaire.ContentsEndpoint, cfg.Vaultaire.Origins, cfg.Indexer.Parallelism)
		IndexerLogger.Debugf("Indexer sleeping for %v seconds.", cfg.Indexer.IndexInterval)
		time.Sleep(time.Duration(cfg.Indexer.IndexInterval) * time.Second)
	}
}
