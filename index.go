package chevalier

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	es "github.com/mattbaird/elastigo/core"
	"strings"
)

type ElasticsearchSource struct {
	Origin string
	Source map[string]string `json:"source"`
}

// GetID returns a (probably) unique ID for an ElasticsearchSource, in
// the form of a sha1 hash of underscore-separated field-value pairs
// separated by newlines.
func (s *ElasticsearchSource) GetID() string {
	tagKeys := make([]string, len(s.Source))
	idx := 0
	for field, value := range s.Source {
		tagKeys[idx] = fmt.Sprintf("%s_%s", field, value)
		idx++
	}
	key := []byte(strings.Join(tagKeys, "\n"))
	hash := sha1.Sum(key)
	id := base64.StdEncoding.EncodeToString(hash[:sha1.Size])
	return id
}

func NewElasticsearchSource(origin string, source *DataSource) *ElasticsearchSource {
	esSource := new(ElasticsearchSource)
	esSource.Origin = origin
	esSource.Source = make(map[string]string, 0)
	for _, tagPtr := range source.Source {
		esSource.Source[*tagPtr.Field] = *tagPtr.Value
	}
	return esSource
}

// Unmarshal turns an ElasticsearchSource (presumably itself unmarshaled
// from a JSON object stored in Elasticsearch) into the equivalent
// DataSource.
func (s *ElasticsearchSource) Unmarshal() *DataSource {
	tags := make([]*DataSource_Tag, len(s.Source))
	idx := 0
	for field, value := range s.Source {
		tags[idx] = NewDataSourceTag(field, value)
		idx++
	}
	pb := NewDataSource(tags)
	return pb
}

func MarshalElasticsearchSources(origin string, b *DataSourceBurst) []*ElasticsearchSource {
	sources := make([]*ElasticsearchSource, len(b.Sources))
	for i, s := range b.Sources {
		esSource := NewElasticsearchSource(origin, s)
		sources[i] = esSource
	}
	return sources
}

// ElasticsearchWriter maintains context for writes to the index.
type ElasticsearchWriter struct {
	indexer   *es.BulkIndexer
	indexName string
	dataType  string
	done      chan bool
}

// NewElasticsearchWriter builds a new Writer. retrySeconds is for the
// bulk indexer. index and dataType can be anything as long as they're
// consistent.
func NewElasticsearchWriter(host string, maxConns int, retrySeconds int, index, dataType string) *ElasticsearchWriter {
	writer := new(ElasticsearchWriter)
	api.Domain = host
	writer.indexer = es.NewBulkIndexerErrors(maxConns, retrySeconds)
	writer.indexName = index
	writer.dataType = dataType
	writer.done = make(chan bool)
	writer.indexer.Run(writer.done)
	return writer
}

// Write queues a DataSource for writing by the bulk indexer.
// Non-blocking.
func (w *ElasticsearchWriter) Write(origin string, source *DataSource) error {
	esSource := NewElasticsearchSource(origin, source)
	err := w.indexer.Index(w.indexName, w.dataType, esSource.GetID(), "", nil, esSource)
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
