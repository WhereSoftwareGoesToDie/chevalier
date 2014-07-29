package chevalier

import (
	"testing"
)

func TestEmptyDataSource(t *testing.T) {
	ds := new(DataSource)
	addr := uint64(42);
	ds.Address = &addr;
	if !ds.Empty() {
		t.Errorf("I think an empty datasource is non-empty.")
	}
	tg_f := "foo"
	tg_v := "bar"
	tg := new(DataSource_Tag)
	tg.Field = &tg_f
	tg.Value = &tg_v
	tgs := make([]*DataSource_Tag, 1)
	tgs[0] = tg
	ds.Source = tgs
	if ds.Empty() {
		t.Errorf("I think a non-empty datasource is empty.")
	}
	ds.Source = ds.Source[:0]
	if !ds.Empty() {
		t.Errorf("I think a cleared datasource is non-empty.")
	}
}

