package chevalier

import (
	"encoding/json"
	"testing"
)

func TestBuildQuery(t *testing.T) {
	engine := NewQueryEngine("localhost", "chevalier_test", "datasource")
	query := new(SourceRequest)
	query.Tags = make([]*SourceRequest_Tag, 2)
	query.Tags[0] = NewSourceRequestTag("hostname", "*.example.com")
	query.Tags[1] = NewSourceRequestTag("metric", "cpu")
	q, err := engine.BuildQuery(query)
	if err != nil {
		t.Errorf("%v", err)
	}
	json, err := json.Marshal(q)
	if err != nil {
		t.Errorf("%v", err)
	}
	expected := `{"from":0,"query":{"bool":{"must":[{"query_string":{"query":"*.example.com","fields":["datasource.hostname"]}},{"query_string":{"query":"cpu","fields":["datasource.metric"]}}]}},"size":0}`
	result := string(json[:])
	if result != expected {
		t.Errorf("Query marshalling mismatch: expected %v, got %v.", expected, result)
	}
}

func TestSanitizeTag(t *testing.T) {
	var tagTests = []struct {
		k string
		v string
		outK string
		outV string
	}{
		{"*", "*", "datasource._all", "*"},
		{"host*", "[test*]^~", "datasource.host*", `\[test*\]\^\~`},
	}
	engine := NewQueryEngine("localhost", "chevalier_test", "datasource")
	failIfInvalid := func(f, v, wantF, wantV string) {
		resF, resV := engine.sanitizeTag(f, v)
		if resF != wantF || resV != wantV {
			t.Errorf("Got %v and %v, wanted %v and %v", resF, resV, wantF, wantV)
		}
	}
	for _, tt := range tagTests {
		failIfInvalid(tt.k, tt.v, tt.outK, tt.outV)
	}
}
