package chevalier

import (
	"errors"
	"encoding/binary"
	"bytes"
	"strings"
	"fmt"
)

type RequestOpCode byte

type ResponseOpCode byte

const (
	ContentsListRequest RequestOpCode = iota
	GenerateNewAddress
	UpdateSourceTag
	RemoveSourceTag
)

const (
	RandomAddress ResponseOpCode = iota
	InvalidContentsOrigin
	ContentsListEntry
	EndOfContentsList
	UpdateSuccess
	RemoveSuccess
)

// ContentsEntry is one address->source-dict mapping.
type ContentsEntry struct {
	address uint64
	tags map[string]string
}

// ContentsResponse is a single message received from the contents daemon.
type ContentsResponse struct {
	opCode ResponseOpCode
	entry *ContentsEntry
}

// unpackTags takes a byteslice representing a comma-separated list of
// field:value pairs (assumed to be UTF-8). It returns the corresponding
// map of fields to values.
func unpackTags(tagData []byte) (map[string]string, error) {
	tags := make(map[string]string, 0)
	if len(tagData) == 0 {
		return tags, nil
	}
	sTags := string(tagData[:])
	pairs := strings.Split(sTags, ",")
	for _, p := range pairs {
		idx := strings.Index(p, ":")
		if idx == -1 || idx+1 == len(p) {
			return tags, errors.New(fmt.Sprintf("Could not parse tag %v.", p))
		}
		field := p[:idx]
		val := p[idx+1:]
		tags[field] = val
	}
	return tags, nil
}

// unpackContentsEntry takes a packet (with the opcode prefix byte
// removed) and returns a ContentsResponse parsed as a
// ContentsListEntry (or an error).
func unpackContentsEntry(packet []byte) (*ContentsResponse,error) {
	buf := bytes.NewBuffer(packet)
	var addr, contentsLen uint64
	err := binary.Read(buf, binary.LittleEndian, &addr)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &contentsLen)
	if err != nil {
		return nil, err
	}
	e := new(ContentsEntry)
	e.address = addr
	e.tags, err = unpackTags(buf.Bytes())
	if err != nil {
		return nil, err
	}
	r := new(ContentsResponse)
	r.opCode = ContentsListEntry
	r.entry  = e
	return r, err
}

// unpackContentsResponse takes a packet received from the contents
// daemon and returns its ContentsResponse representation or an error.
// Currently errors on packet types other than ContentsListEntry and
// EndOfContentsList as those are the only ones we should be seeing.
func unpackContentsResponse(packet []byte) (*ContentsResponse,error) {
	if len(packet) == 0 {
		return nil, errors.New("Empty packet.")
	}
	opcode := ResponseOpCode(packet[0])
	typeError := errors.New("Unexpected packet type - expecting ContentsListEntry or EndOfContentsList.")
	switch (opcode) {
	case RandomAddress:
		return nil, typeError
	case InvalidContentsOrigin:
		return nil, typeError
	case ContentsListEntry:
		return unpackContentsEntry(packet[1:])
	case EndOfContentsList:
		res := new(ContentsResponse)
		res.opCode = EndOfContentsList
		return res, nil
	case UpdateSuccess:
		return nil, typeError
	case RemoveSuccess:
		return nil, typeError
	default:
		return nil, errors.New("Invalid response opcode.")
	}
}
