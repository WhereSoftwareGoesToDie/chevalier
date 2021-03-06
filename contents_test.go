package chevalier

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDecodeContentsEnd(t *testing.T) {
	pkt := make([]byte, 1)
	pkt[0] = 3
	res, err := unpackContentsResponse(pkt)
	if err != nil {
		t.Errorf("Error decoding packet: %v", err)
	} else if res.opCode != EndOfContentsList {
		t.Errorf("Wrong opcode: expected EndOfContentsList, got %v", res.opCode)
	}
}

func TestDecodeContentsEntry(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.WriteByte(2) // ContentsListEntry
	addr := int64(42)
	dict := []byte("foo:bar,baz:quux,")
	dictLen := int64(len(dict))
	binary.Write(buf, binary.LittleEndian, addr)
	binary.Write(buf, binary.LittleEndian, dictLen)
	buf.Write(dict)
	res, err := unpackContentsResponse(buf.Bytes())
	if err != nil {
		t.Errorf("Error decoding packet: %v", err)
	} else if res.opCode != ContentsListEntry {
		t.Errorf("Wrong opcode: expected ContentsListEntry, got %v", res.opCode)
	}
	if res.entry == nil {
		t.Errorf("Error parsing contents entry, got a nil pointer")
	}
	t.Logf("tags: %v", res.entry.tags)
	if v, ok := res.entry.tags["foo"]; ok {
		if v != "bar" {
			t.Errorf("Corrupted tag: expected bar, got %v (in source dict %v)", v, res.entry.tags)
		}
	} else {
		t.Errorf("Expected tag foo in source dict %v", res.entry.tags)
	}
}

func TestDecodeEmptyTag(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.WriteByte(2) // ContentsListEntry
	addr := int64(42)
	dict := []byte("foo:,")
	dictLen := int64(len(dict))
	binary.Write(buf, binary.LittleEndian, addr)
	binary.Write(buf, binary.LittleEndian, dictLen)
	buf.Write(dict)
	res, err := unpackContentsResponse(buf.Bytes())
	if err != nil {
		t.Errorf("Error decoding packet: %v", err)
	}
	if v, ok := res.entry.tags["foo"]; ok {
		if v != "" {
			t.Errorf("Corrupted tag: expected empty string, got %v (in source dict %v)", v, res.entry.tags)
		}
	} else {
		t.Errorf("Empty value failure: expected tag foo in source dict %v", res.entry.tags)
	}
}
