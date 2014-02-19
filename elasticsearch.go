package chevalier

import (
	es "github.com/mattbaird/elastigo/core"
	"fmt"
	"strings"
	"crypto/sha1"
)

type ElasticsearchSourceTag struct {
	Field string
	Value string
}

type ElasticsearchSource struct {
	Source []ElasticsearchSourceTag
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
	return string(hash[:sha1.Size])
}

func NewElasticsearchSource(source *DataSource) *ElasticsearchSource {
	esSource := new(ElasticsearchSource)
	esSource.Source = make([]ElasticsearchSourceTag, len(source.Source))
	for i, tagPtr := range source.Source {
		esSource.Source[i].Field = *tagPtr.Field
		esSource.Source[i].Value = *tagPtr.Field
	}
	return esSource
}

type ElasticsearchWriter struct {
	indexer *es.BulkIndexer
	indexName string
	dataType string
}

func NewElasticsearchWriter(host string, maxConns int, retrySeconds int, index, dataType string) *ElasticsearchWriter {
	writer := new(ElasticsearchWriter)
	writer.indexer = es.NewBulkIndexerErrors(maxConns, retrySeconds)
	writer.indexName = index
	writer.dataType = dataType
	return writer
}

func (w *ElasticsearchWriter) Write(source *DataSource) {
	esSource := NewElasticsearchSource(source)
	w.indexer.Index(w.indexName, w.dataType, esSource.GetID(), "", nil, esSource)
}
