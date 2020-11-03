package udp

import (
	"testing"
)

func TestChecksum(t *testing.T) {
	packet := NewPacket([]byte("hello world"), 57090, 8000, 3232235623, 3232235521)
	if packet.header.Checksum != 60925 {
		t.Error("invalid checksum: ", packet.header.Checksum)
	}
}
