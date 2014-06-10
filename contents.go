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

type ContentsEntry struct {
	address uint64
	tags map[string]string
}

type ContentsResponse struct {
	opCode byte
	entry *ContentsEntry
}

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
	r.opCode = byte(ContentsListEntry)
	r.entry  = e
	return r, err
}

func unpackContentsResponse(packet []byte) (*ContentsResponse,error) {
	if len(packet) == 0 {
		return nil, errors.New("Empty packet.")
	}
	opcode := ResponseOpCode(packet[0])
	switch (opcode) {
	case RandomAddress:
		return nil, nil
	case InvalidContentsOrigin:
		return nil, nil
	case ContentsListEntry:
		return unpackContentsEntry(packet[1:])
	case EndOfContentsList:
		return nil, nil
	case UpdateSuccess:
		return nil, nil
	case RemoveSuccess:
		return nil, nil
	default:
		return nil, errors.New("Invalid response opcode.")
	}
}
