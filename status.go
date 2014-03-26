package chevalier

import (
	"encoding/json"
	"fmt"
	"time"
	"errors"

	"code.google.com/p/goprotobuf/proto"
	es "github.com/mattbaird/elastigo/core"
)

func MarshalStatusResponse(s *StatusResponse) ([]byte, error) {
	return proto.Marshal(s)
}

func (s *StatusResponse) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

func UnmarshalStatusResponse(b []byte) (*StatusResponse, error) {
	s := new(StatusResponse)
	err := proto.Unmarshal(b, s)
	return s, err
}

func NewStatusResponse() (*StatusResponse) {
	s := new(StatusResponse)
	s.Origins = make([]*StatusResponse_Origin, 0)
	s.Errors = make([]string, 0)
	return s
}

func (s *StatusResponse) addError(err error) {
	s.Errors = append(s.Errors, fmt.Sprintf("%v", err))
}

func NewStatusResponse_Origin(origin string, sources uint64, updated time.Time) *StatusResponse_Origin {
	o := new(StatusResponse_Origin)
	originR := origin
	sourcesR := sources
	updatedR := uint64(updated.UnixNano())
	o.Origin = &originR
	o.Sources = &sourcesR
	o.LastUpdated = &updatedR
	return o
}

// runOriginQuery returns an elastigo/core.SearchResult for the
// specified origin in the metadata index.
func (e *QueryEngine) runOriginQuery(origin string) (*es.SearchResult, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"query_string": map[string]interface{} {
						"default_field" : "origin",
						"query" : origin,
					},
				},
			},
		},
	}
	var args map[string]interface{}
	r, err := es.SearchRequest(e.metaIndex, e.originType, args, query)
	return &r, err
}

// GetOriginMetadata returns an ElasticsearchOrigin object for the
// specified origin. 
func (e *QueryEngine) GetOriginMetadata(origin string) (*ElasticsearchOrigin, error) {
	res, err := e.runOriginQuery(origin)
	if err != nil {
		return nil, err
	}
	if len(res.Hits.Hits) == 0 {
		errMsg := fmt.Sprintf("no metadata available for origin %v", origin)
		return nil, errors.New(errMsg)
	}
	mRes := res.Hits.Hits[0]
	om := new(ElasticsearchOrigin)
	err = json.Unmarshal(*mRes.Source, om)
	if err != nil {
		return nil, err
	}
	return om, nil
}

func (e *QueryEngine) GetStatus(origins []string) (*StatusResponse) {
	s := new(StatusResponse)
	for _, o := range origins {
		esOrigin, err := e.GetOriginMetadata(o)
		if err != nil {
			s.addError(err)
			continue
		}
		pbOrigin := NewStatusResponse_Origin(esOrigin.Origin, esOrigin.Count, esOrigin.LastUpdated)
		s.Origins = append(s.Origins, pbOrigin)
	}
	return s
}
