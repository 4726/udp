package udp

import (
	"encoding/binary"
)

const UDPProtocolType = 17

type Packet struct {
	header  Header
	pHeader PseudoHeader
}

type Header struct {
	SourcePort, DestinationPort uint16
	Length, Checksum            uint16
	Data                        []byte
}

type PseudoHeader struct {
	SourceAddr, DestinationAddr uint32
	Length                      uint16
}

func NewPacket(data []byte, sourcePort, destPort uint16, sourceAddr, destAddr uint32) Packet {
	header := Header{
		SourcePort:      sourcePort,
		DestinationPort: destPort,
		Data:            data,
	}
	pHeader := PseudoHeader{
		SourceAddr:      sourceAddr,
		DestinationAddr: destAddr,
	}
	header.Length = uint16(2 + 2 + 2 + 2 + len(data))
	pHeader.Length = header.Length
	header.Checksum = checksum(header, pHeader)

	return Packet{
		header:  header,
		pHeader: pHeader,
	}
}

func (p Packet) Pack() []byte {
	res := make([]byte, 2+2+2+2+len(p.header.Data))

	binary.BigEndian.PutUint16(res, p.header.SourcePort)
	binary.BigEndian.PutUint16(res[2:], p.header.DestinationPort)
	binary.BigEndian.PutUint16(res[4:], p.header.Length)
	binary.BigEndian.PutUint16(res[6:], p.header.Checksum)
	copy(res[8:], p.header.Data)

	return res
}

func checksum(header Header, pHeader PseudoHeader) uint16 {
	sourceAddr1, sourceAddr2 := uint32ToTwoUint16(pHeader.SourceAddr)
	destinationAddr1, destinationAddr2 := uint32ToTwoUint16(pHeader.DestinationAddr)
	total := uint32(sourceAddr1) +
		uint32(sourceAddr2) +
		uint32(destinationAddr1) +
		uint32(destinationAddr2) +
		uint32(UDPProtocolType) +
		uint32(pHeader.Length) +
		uint32(header.SourcePort) +
		uint32(header.DestinationPort) +
		uint32(header.Length)

	for i := 0; i < len(header.Data); i += 2 {
		if i+1 == len(header.Data) {
			b := []byte{header.Data[len(header.Data)-1], 0}
			total += uint32(binary.BigEndian.Uint16(b))
		} else {
			total += uint32(binary.BigEndian.Uint16(header.Data[i : i+2]))
		}
	}

	for total >= 1<<16 {
		num1, num2 := uint32ToTwoUint16(total)
		total = uint32(num1) + uint32(num2)
	}

	return uint16(total) ^ ((1 << 16) - 1)
}

func uint32ToTwoUint16(n uint32) (uint16, uint16) {
	var num1, num2 uint16

	for i := 0; i < 16; i++ {
		if (n>>i)%2 == 1 {
			num1 |= (1 << i)
		}
	}

	for i := 16; i < 32; i++ {
		if (n>>i)%2 == 1 {
			num2 |= (1 << (i - 16))
		}
	}

	return num1, num2
}
