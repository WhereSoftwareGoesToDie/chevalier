package chevalier

import (
	"github.com/mattbaird/elastigo/api"
	es "github.com/mattbaird/elastigo/core"
	"github.com/mattbaird/elastigo/search"
)

type QueryEngine struct {
	indexName string
	dataType  string
}

func NewQueryEngine(host, indexName, dataType string) *QueryEngine {
	e := new(QueryEngine)
	e.indexName = indexName
	e.dataType = dataType
	api.Domain = host
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
		"size": 100000,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": tagQueries,
			},
		},
	}
	return SourceQuery(query)
}

func (e *QueryEngine) RunSourceRequest(req *SourceRequest) (*es.SearchResult, error) {
	q := e.BuildQuery(req)
	res, err := es.SearchRequest(false, e.indexName, e.dataType, q, "", 0)
	return &res, err
}

func FmtResult(result *es.SearchResult) []string {
	results := make([]string, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		results[i] = string(hit.Source[:])
	}
	return results
}
