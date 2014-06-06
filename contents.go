package chevalier

import (
	"errors"
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
	sources map[string]string
}

type ContentsResponse struct {
	opCode byte
	entries []ContentsEntry
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
		return nil, nil
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
