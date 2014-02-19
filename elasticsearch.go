package chevalier

import (
	es "github.com/mattbaird/elastigo/core"
	"github.com/mattbaird/elastigo/api"
	"fmt"
	"strings"
	"encoding/base64"
	"crypto/sha1"
)

type ElasticsearchSourceTag struct {
	Field string `json:"tag_field"`
	Value string `json:"tag_value"`
}

type ElasticsearchSource struct {
	Source []ElasticsearchSourceTag `json:"source"`
}

// GetID returns a (probably) unique ID for an ElasticsearchSource, in
// the form of a sha1 hash of underscore-separated field-value pairs
// separated by newlines.
func (s *ElasticsearchSource) GetID() string {
	tagKeys := make([]string, len(s.Source))
	for i, tag := range s.Source {
		tagKeys[i] = fmt.Sprintf("%s_%s", tag.Field, tag.Value)
	}
	key := []byte(strings.Join(tagKeys, "\n"))
	hash := sha1.Sum(key)
	id := base64.StdEncoding.EncodeToString(hash[:sha1.Size])
	return id
}

func NewElasticsearchSource(source *DataSource) *ElasticsearchSource {
	esSource := new(ElasticsearchSource)
	esSource.Source = make([]ElasticsearchSourceTag, len(source.Source))
	for i, tagPtr := range source.Source {
		esSource.Source[i].Field = *tagPtr.Field
		esSource.Source[i].Value = *tagPtr.Value
	}
	return esSource
}

type ElasticsearchWriter struct {
	indexer *es.BulkIndexer
	indexName string
	dataType string
	done chan bool
}

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

func (w *ElasticsearchWriter) Write(source *DataSource) error {
	esSource := NewElasticsearchSource(source)
	err := w.indexer.Index(w.indexName, w.dataType, esSource.GetID(), "", nil, esSource)
	return err
}

func (w *ElasticsearchWriter) WaitDone() {
	<-w.done
}
