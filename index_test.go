package chevalier

import (
	"testing"
)

func TestGetID(t *testing.T) {
	tags := make([]*DataSource_Tag, 2)
	tags[0] = NewDataSourceTag("foo", "bar")
	tags[1] = NewDataSourceTag("baz", "quux")
	source := NewDataSource(tags)
	esSource := NewElasticsearchSource("ABCDEF", source)
	result := esSource.GetID()
	expected := "0hhETFpfnLKUZTgfDJDJ1Sq9GdY="
	if result != expected {
		t.Errorf("Got ID %v (expected %v)", result, expected)
	}
}
