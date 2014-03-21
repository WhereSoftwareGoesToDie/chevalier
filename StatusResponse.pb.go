// Code generated by protoc-gen-go.
// source: StatusResponse.proto
// DO NOT EDIT!

package chevalier

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

// Sent as a response to a status request.
type StatusResponse struct {
	// All origins currently in the index.
	Origins          []*StatusResponse_Origin `protobuf:"bytes,1,rep,name=origins" json:"origins,omitempty"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (m *StatusResponse) Reset()         { *m = StatusResponse{} }
func (m *StatusResponse) String() string { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()    {}

func (m *StatusResponse) GetOrigins() []*StatusResponse_Origin {
	if m != nil {
		return m.Origins
	}
	return nil
}

type StatusResponse_Origin struct {
	// Origin name.
	Origin *string `protobuf:"bytes,1,req,name=origin" json:"origin,omitempty"`
	// Number of sources for this origin.
	Sources *int64 `protobuf:"varint,2,req,name=sources" json:"sources,omitempty"`
	// Nanosecond-precision timestamp of last update.
	LastUpdated      *uint64 `protobuf:"fixed64,3,opt,name=last_updated" json:"last_updated,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *StatusResponse_Origin) Reset()         { *m = StatusResponse_Origin{} }
func (m *StatusResponse_Origin) String() string { return proto.CompactTextString(m) }
func (*StatusResponse_Origin) ProtoMessage()    {}

func (m *StatusResponse_Origin) GetOrigin() string {
	if m != nil && m.Origin != nil {
		return *m.Origin
	}
	return ""
}

func (m *StatusResponse_Origin) GetSources() int64 {
	if m != nil && m.Sources != nil {
		return *m.Sources
	}
	return 0
}

func (m *StatusResponse_Origin) GetLastUpdated() uint64 {
	if m != nil && m.LastUpdated != nil {
		return *m.LastUpdated
	}
	return 0
}

func init() {
}