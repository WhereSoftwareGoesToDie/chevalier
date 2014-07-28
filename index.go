package chevalier

import (
	"fmt"
	"strconv"
	"time"

	"github.com/anchor/elastigo/api"
	es "github.com/anchor/elastigo/core"
)

// ElasticsearchSource is the type used to serialize sources for
// indexing.
type ElasticsearchSource struct {
	Origin string
	// Address in Vaultaire.
	Address string
	Source  map[string]string `json:"source"`
}

// ElasticsearchOrigin stores metadata for each origin.
type ElasticsearchOrigin struct {
	Origin      string    `json:"origin"`
	Count       uint64    `json:"count"`
	Address     uint64    `json:"address"`
	LastUpdated time.Time `json:"last_updated"`
}

func NewElasticsearchOrigin(origin string, count uint64, updated time.Time) *ElasticsearchOrigin {
	o := new(ElasticsearchOrigin)
	o.Origin = origin
	o.Count = count
	o.LastUpdated = updated
	return o
}

// GetID returns a (probably) unique ID for an ElasticsearchSource, in
// the form of a sha1 hash of underscore-separated field-value pairs
// separated by newlines.
func (s *ElasticsearchSource) GetID() string {
	k := fmt.Sprintf("%v/%v", s.Origin, s.Address)
	return k
}

// Unmarshal turns an ElasticsearchSource (presumably itself unmarshaled
// from a JSON object stored in Elasticsearch) into the equivalent
// DataSource.
func (s *ElasticsearchSource) Unmarshal() (*DataSource, error) {
	tags := make([]*DataSource_Tag, len(s.Source))
	idx := 0
	for field, value := range s.Source {
		tags[idx] = NewDataSourceTag(field, value)
		idx++
	}
	pb := NewDataSource(tags)
	addr, err := strconv.ParseUint(s.Address, 10, 64)
	if err != nil {
		return nil, err
	}
	pb.Address = &addr
	return pb, nil
}

// ElasticsearchWriter maintains context for writes to the index.
type ElasticsearchWriter struct {
	indexer   *es.BulkIndexer
	indexName string
	// Metadata index
	metaIndex  string
	dataType   string
	originType string
	done       chan bool
}

// NewElasticsearchWriter builds a new Writer. retrySeconds is for the
// bulk indexer. index and dataType can be anything as long as they're
// consistent.
func NewElasticsearchWriter(host string, maxConns int, retrySeconds int, index, metaIndex, dataType string) *ElasticsearchWriter {
	writer := new(ElasticsearchWriter)
	api.Domain = host
	writer.indexer = es.NewBulkIndexerErrors(maxConns, retrySeconds)
	writer.indexName = index
	writer.metaIndex = metaIndex
	writer.dataType = dataType
	writer.originType = "chevalier_origin"
	writer.done = make(chan bool)
	writer.indexer.Run(writer.done)
	return writer
}

func (w *ElasticsearchWriter) UpdateOrigin(origin string, count uint64) error {
	o := NewElasticsearchOrigin(origin, count, time.Now())
	update := map[string]interface{}{
		"doc":           o,
		"doc_as_upsert": true,
	}
	err := w.indexer.Update(w.metaIndex, w.originType, origin, "", nil, update, true)
	return err
}

// Write queues a DataSource for writing by the bulk indexer.
// Non-blocking.
func (w *ElasticsearchWriter) Write(origin string, source *ElasticsearchSource) error {
	update := map[string]interface{}{
		"doc":           source,
		"doc_as_upsert": true,
	}
	err := w.indexer.Update(w.indexName, w.dataType, source.GetID(), "", nil, update, true)
	return err
}

// Shutdown signals the bulk indexer to flush all pending writes.
func (w *ElasticsearchWriter) Shutdown() {
	w.done <- true
}

// GetErrorChan returns the channel the bulk indexer writes errors to.
func (w *ElasticsearchWriter) GetErrorChan() chan *es.ErrorBuffer {
	return w.indexer.ErrorChannel
}
