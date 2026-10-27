package fixture

import (
	"net"

	impl "github.com/grandper/go-screenplay-udp/internal/ability"
	"github.com/grandper/go-screenplay-udp/udp"
)

// CreateFakeDatagram creates a fake datagram holding the text, sent from the
// address. An address that cannot be resolved leaves the sender unknown.
func CreateFakeDatagram(text, sender string) *udp.Datagram {
	return CreateFakeBinaryDatagram([]byte(text), sender)
}

// CreateFakeBinaryDatagram creates a fake datagram holding the bytes, sent
// from the address. An address that cannot be resolved leaves the sender
// unknown.
func CreateFakeBinaryDatagram(data []byte, sender string) *udp.Datagram {
	address, err := net.ResolveUDPAddr("udp", sender)
	if err != nil {
		address = nil
	}

	return impl.NewDatagram(data, address)
}
