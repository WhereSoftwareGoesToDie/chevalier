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
