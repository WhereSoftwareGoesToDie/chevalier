package chevalier

import (
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
