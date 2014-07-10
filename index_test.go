package chevalier

import (
	"testing"
)

func TestGetID(t *testing.T) {
	source := new(ElasticsearchSource)
	source.Source = make(map[string]string, 0)
	source.Origin = "ABCDEF"
	source.Address = "42"
	source.Source["foo"] = "bar"
	source.Source["baz"] = "quux"
	result := source.GetID()
	expected := "ABCDEF/42"
	if result != expected {
		t.Errorf("Got ID %v (expected %v)", result, expected)
	}
}
