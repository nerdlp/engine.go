package engine

import (
	"github.com/gammazero/deque"
)

type transportClient interface {
	// get packet from client
	getPacket() <-chan *packet
	// receivePacket should be invoke to notify that a new packet received from client.
	receivePacket(packet *packet) error
	// send packet to the client
	sendPacket(packet *packet) error
}

type pollingClient struct {
	buffers           deque.Deque[*packet]
	maxPayload        int
	receivePacketChan chan *packet
}

func (c *pollingClient) getPacket() <-chan *packet {
	return c.receivePacketChan
}

func (c *pollingClient) receivePacket(packet *packet) error {
	c.receivePacketChan <- packet
	return nil
}

// For polling transport, sendPacket does not send packet to client, it just adds to the buffers and wait for the client to poll it
func (c *pollingClient) sendPacket(packet *packet) error {
	c.buffers.PushBack(packet)
	return nil
}

// compress packet for better response
func (c *pollingClient) compressPacket() []byte {
	// compress packet
	data := make([]byte, 0, c.maxPayload)

	for c.buffers.Len() != 0 {
		nextPacket := c.buffers.Front().encodePolling()
		if len(nextPacket)+len(data)+1 > c.maxPayload {
			break
		}
		data = append(data, recordSeperator)
		data = append(data, nextPacket...)
		c.buffers.PopFront()
	}
	return data
}

// decompress packet for better response
func (c *pollingClient) decompressPacket(data []byte) []*packet {
	start := 0
	packets := make([]*packet, 0)
	for i := range data {
		if data[i] == recordSeperator {
			packet := new(packet)
			packet.decodePolling(data[start:i])
			start = i + 1
			packets = append(packets, packet)
		}
	}
	return packets
}
