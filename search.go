package chevalier

import (
	"encoding/json"
	"github.com/mattbaird/elastigo/search"
	"github.com/mattbaird/elastigo/api"
)

func NewSourceRequestTag(field, value string) *SourceRequest_Tag {
	tag := new(SourceRequest_Tag)
	f := field
	v := value
	tag.Field = &f
	tag.Value = &v
	return tag
}

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

func (e *QueryEngine) BuildQuery(req *SourceRequest) ([]byte, error) {
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
	data, err := json.Marshal(query)
	return data, err
}
