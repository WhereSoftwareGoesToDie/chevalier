package chevalier

import (
	"code.google.com/p/goprotobuf/proto"
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
