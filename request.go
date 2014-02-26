package chevalier

import (
	"code.google.com/p/goprotobuf/proto"
)

// NewSourceRequestTag builds a SourceRequest_Tag from a field and a
// value.
func NewSourceRequestTag(field, value string) *SourceRequest_Tag {
	tag := new(SourceRequest_Tag)
	f := field
	v := value
	tag.Field = &f
	tag.Value = &v
	return tag
}

// NewSourceRequest builds a SourceRequest from a slice of tags.
func NewSourceRequest(tags []*SourceRequest_Tag) *SourceRequest {
	req := new(SourceRequest)
	req.Tags = tags
	return req
}

func MarshalSourceRequest(req *SourceRequest) ([]byte, error) {
	marshalled, err := proto.Marshal(req)
	return marshalled, err
}

func UnmarshalSourceRequest(packet []byte) (*SourceRequest, error) {
	source := new(SourceRequest)
	err := proto.Unmarshal(packet, source)
	return source, err
}
