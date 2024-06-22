package engine

type packet struct {
	ptype packetType
	data  []byte
}

// return new packet
func newPacket(ptype packetType, data []byte) *packet {
	return &packet{
		ptype: ptype,
		data:  data,
	}
}

// encoding packet for polling transport
func (p *packet) encodePolling() []byte {
	encodedPacket := make([]byte, 0, 1+len(p.data))
	encodedPacket = append(encodedPacket, p.ptype)
	encodedPacket = append(encodedPacket, p.data...)
	return encodedPacket
}

// encoding packet for polling transport
func (p *packet) decodePolling(data []byte) {
	p.ptype = data[0]
	p.data = data[1:]
}

// encoding packet for websocket transport
func (p *packet) encodeWebSocket() []byte {
	encodedPacket := make([]byte, 0, 1+len(p.data))
	encodedPacket = append(encodedPacket, p.ptype)
	encodedPacket = append(encodedPacket, p.data...)
	return encodedPacket
}
