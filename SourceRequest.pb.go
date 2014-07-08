// Code generated by protoc-gen-go.
// source: SourceRequest.proto
// DO NOT EDIT!

package chevalier

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

// This message is sent by Chevalier clients; a DataSourceBurst is sent
// in response.
type SourceRequest struct {
	// Tags to use in search (as an 'and' query). If `query_string`
	// is specified, the content of this field is ignored.
	Tags []*SourceRequest_Tag `protobuf:"bytes,1,rep,name=tags" json:"tags,omitempty"`
	// Page to return results from. If not specified, 0 is assumed.
	StartPage *int64 `protobuf:"varint,2,opt,name=start_page" json:"start_page,omitempty"`
	// Page to return results from. If not specified, all results
	// are returned in one page.
	SourcesPerPage *int64 `protobuf:"varint,3,opt,name=sources_per_page" json:"sources_per_page,omitempty"`
	// Elasticsearch query string to use. If specified, the content
	// of `tags` will be ignored.
	QueryString *string `protobuf:"bytes,5,opt,name=query_string" json:"query_string,omitempty"`
	// Vaultaire address to look up. If specified, the `tags` and
	// `query_string` fields will be ignored.
	Address          *uint64 `protobuf:"fixed64,6,opt,name=address" json:"address,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SourceRequest) Reset()         { *m = SourceRequest{} }
func (m *SourceRequest) String() string { return proto.CompactTextString(m) }
func (*SourceRequest) ProtoMessage()    {}

func (m *SourceRequest) GetTags() []*SourceRequest_Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *SourceRequest) GetStartPage() int64 {
	if m != nil && m.StartPage != nil {
		return *m.StartPage
	}
	return 0
}

func (m *SourceRequest) GetSourcesPerPage() int64 {
	if m != nil && m.SourcesPerPage != nil {
		return *m.SourcesPerPage
	}
	return 0
}

func (m *SourceRequest) GetQueryString() string {
	if m != nil && m.QueryString != nil {
		return *m.QueryString
	}
	return ""
}

func (m *SourceRequest) GetAddress() uint64 {
	if m != nil && m.Address != nil {
		return *m.Address
	}
	return 0
}

type SourceRequest_Tag struct {
	Field            *string `protobuf:"bytes,1,req,name=field" json:"field,omitempty"`
	Value            *string `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SourceRequest_Tag) Reset()         { *m = SourceRequest_Tag{} }
func (m *SourceRequest_Tag) String() string { return proto.CompactTextString(m) }
func (*SourceRequest_Tag) ProtoMessage()    {}

func (m *SourceRequest_Tag) GetField() string {
	if m != nil && m.Field != nil {
		return *m.Field
	}
	return ""
}

func (m *SourceRequest_Tag) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func init() {
}
