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
	json, err := json.Marshal(engine.BuildQuery(query))
	if err != nil {
		t.Errorf("%v", err)
	}
	expected := `{"query":{"bool":{"must":[{"query_string":{"query":"*.example.com","fields":["hostname"]}},{"query_string":{"query":"cpu","fields":["metric"]}}]}}}`
	result := string(json[:])
	if result != expected {
		t.Errorf("Query marshalling mismatch: expected %v, got %v.", expected, result)
	}
}
