// Code generated by protoc-gen-go.
// source: DataSource.proto
// DO NOT EDIT!

package chevalier

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

// Returned as a response to the chevalier client.
type DataSourceBurst struct {
	// All sources that matched the query received (paginated
	// according to `start_page` and `sources_per_page` if they are
	// set in the request).
	Sources []*DataSource `protobuf:"bytes,1,rep,name=sources" json:"sources,omitempty"`
	// Error message - if present, some aspect of the request
	// failed.
	Error            *string `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *DataSourceBurst) Reset()         { *m = DataSourceBurst{} }
func (m *DataSourceBurst) String() string { return proto.CompactTextString(m) }
func (*DataSourceBurst) ProtoMessage()    {}

func (m *DataSourceBurst) GetSources() []*DataSource {
	if m != nil {
		return m.Sources
	}
	return nil
}

func (m *DataSourceBurst) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

type DataSource struct {
	// Source tags. There can be an arbitrary number of these.
	// Tags which affect presentation rather than identity should be
	// underscore-prefixed.
	Source []*DataSource_Tag `protobuf:"bytes,1,rep,name=source" json:"source,omitempty"`
	// Unique identifier for this data source within Vaultaire.
	Address          *uint64 `protobuf:"fixed64,3,req,name=address" json:"address,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *DataSource) Reset()         { *m = DataSource{} }
func (m *DataSource) String() string { return proto.CompactTextString(m) }
func (*DataSource) ProtoMessage()    {}

func (m *DataSource) GetSource() []*DataSource_Tag {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *DataSource) GetAddress() uint64 {
	if m != nil && m.Address != nil {
		return *m.Address
	}
	return 0
}

type DataSource_Tag struct {
	Field            *string `protobuf:"bytes,1,req,name=field" json:"field,omitempty"`
	Value            *string `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *DataSource_Tag) Reset()         { *m = DataSource_Tag{} }
func (m *DataSource_Tag) String() string { return proto.CompactTextString(m) }
func (*DataSource_Tag) ProtoMessage()    {}

func (m *DataSource_Tag) GetField() string {
	if m != nil && m.Field != nil {
		return *m.Field
	}
	return ""
}

func (m *DataSource_Tag) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func init() {
}
