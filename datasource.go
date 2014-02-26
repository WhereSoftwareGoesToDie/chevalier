package chevalier

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/anchor/bletchley/dataframe"
)

func UnmarshalSource(packet []byte) (*DataSource, error) {
	source := new(DataSource)
	err := proto.Unmarshal(packet, source)
	return source, err
}

func UnmarshalSourceBurst(packet []byte) (*DataSourceBurst, error) {
	burst := new(DataSourceBurst)
	err := proto.Unmarshal(packet, burst)
	return burst, err
}

func BuildSourceBurst(sources []*DataSource) *DataSourceBurst {
	burst := new(DataSourceBurst)
	burst.Sources = sources
	return burst
}

func MarshalSourceBurst(burst *DataSourceBurst) ([]byte, error) {
	marshalled, err := proto.Marshal(burst)
	return marshalled, err
}

// DataFrameSource takes a Vaultaire DataFrame and returns just its
// source component as a DataSource.
func DataFrameSource(frame *dataframe.DataFrame) *DataSource {
	source := new(DataSource)
	source.Source = make([]*DataSource_Tag, len(frame.Source))
	for i, tag := range frame.Source {
		source.Source[i] = new(DataSource_Tag)
		source.Source[i].Field = tag.Field
		source.Source[i].Value = tag.Value
	}
	return source
}

func NewDataSourceTag(field, value string) *DataSource_Tag {
	tag := new(DataSource_Tag)
	f := field
	v := value
	tag.Field = &f
	tag.Value = &v
	return tag
}

func NewDataSource(tags []*DataSource_Tag) *DataSource {
	source := new(DataSource)
	source.Source = tags
	return source
}
