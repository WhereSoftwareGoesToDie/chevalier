package chevalier

import (
	"encoding/json"
	"testing"
)

func TestBuildQuery(t *testing.T) {
	engine := NewQueryEngine("localhost", "chevalier_test", "datasource", "chevalier_metadata")
	query := new(SourceRequest)
	query.Tags = make([]*SourceRequest_Tag, 2)
	query.Tags[0] = NewSourceRequestTag("hostname", "*.example.com")
	query.Tags[1] = NewSourceRequestTag("metric", "cpu")
	q, err := engine.BuildQuery("ABCDEF", query)
	if err != nil {
		t.Errorf("%v", err)
	}
	json, err := json.Marshal(q)
	if err != nil {
		t.Errorf("%v", err)
	}
	expected := `{"from":0,"query":{"bool":{"must":[{"query_string":{"analyzer":"keyword","fields":["datasource.hostname"],"query":"*.example.com"}},{"query_string":{"analyzer":"keyword","fields":["datasource.metric"],"query":"cpu"}},{"query_string":{"fields":["Origin"],"query":"ABCDEF"}}]}},"size":0}`
	result := string(json[:])
	if result != expected {
		t.Errorf("Query marshalling mismatch: expected %v, got %v.", expected, result)
	}
}

func TestSanitizeTag(t *testing.T) {
	var tagTests = []struct {
		k    string
		v    string
		outK string
		outV string
	}{
		{"*", "*", "datasource._all", "*"},
		{"host*", "[test*]^~", "datasource.host*", `\[test*\]\^\~`},
	}
	engine := NewQueryEngine("localhost", "chevalier_test", "datasource", "chevalier_metadata")
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
