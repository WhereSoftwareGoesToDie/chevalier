package chevalier

import (
	"encoding/json"
	"github.com/mattbaird/elastigo/api"
	es "github.com/mattbaird/elastigo/core"
	"github.com/mattbaird/elastigo/search"
	"log"
	"time"
)

type QueryEngine struct {
	indexName      string
	dataType       string
	nSources       int
	updateInterval time.Duration
}

func NewQueryEngine(host, indexName, dataType string) *QueryEngine {
	e := new(QueryEngine)
	e.indexName = indexName
	e.dataType = dataType
	api.Domain = host
	e.updateSourceCount()
	e.updateInterval = time.Second * 10
	go e.updateForever()
	return e
}

func (e *QueryEngine) buildTagQuery(tag *SourceRequest_Tag) *search.QueryDsl {
	qs := new(search.QueryString)
	qs.Fields = make([]string, 0)
	qs.Fields = append(qs.Fields, *tag.Field)
	qs.Query = *tag.Value
	q := search.Query().Qs(qs)
	return q
}

type SourceQuery map[string]interface{}

func (e *QueryEngine) BuildQuery(req *SourceRequest) SourceQuery {
	_ = search.Search(e.indexName).Type(e.dataType)
	tags := req.GetTags()
	tagQueries := make([]*search.QueryDsl, len(tags))
	for i, tag := range tags {
		tagQueries[i] = e.buildTagQuery(tag)
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": tagQueries,
			},
		},
	}
	return SourceQuery(query)
}

func (e *QueryEngine) runSourceRequest(req *SourceRequest) (*es.SearchResult, error) {
	q := e.BuildQuery(req)
	res, err := es.SearchRequest(false, e.indexName, e.dataType, q, "", 0)
	return &res, err
}

func (e *QueryEngine) GetSources(req *SourceRequest) (*DataSourceBurst, error) {
	res, err := e.runSourceRequest(req)
	if err != nil {
		return nil, err
	}
	sources := make([]*DataSource, len(res.Hits.Hits))
	for i, hit := range res.Hits.Hits {
		sources[i] = new(DataSource)
		err = json.Unmarshal(hit.Source, sources[i])
		if err != nil {
			return nil, err
		}
	}
	burst := BuildSourceBurst(sources)
	return burst, nil
}

func FmtResult(result *es.SearchResult) []string {
	results := make([]string, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		results[i] = string(hit.Source[:])
	}
	return results
}

func (e *QueryEngine) updateSourceCount() error {
	resp, err := es.Count(false, e.indexName, e.dataType)
	e.nSources = resp.Count
	return err
}

func (e *QueryEngine) updateForever() {
	for true {
		time.Sleep(e.updateInterval)
		err := e.updateSourceCount()
		if err != nil {
			log.Printf("Error updating source count: %v", err)
		}
	}
}
