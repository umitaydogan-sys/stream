package ingest

import (
	"net"

	"github.com/fluxstream/fluxstream/internal/media"
)

// StreamHandler defines the interface for handling stream events from any ingest protocol
type StreamHandler interface {
	OnPublish(streamKey string, conn net.Conn) error
	OnUnpublish(streamKey string)
	OnPacket(streamKey string, pkt *media.Packet)
}
