package fixture

import (
	"net"

	"github.com/grandper/go-screenplay-udp/ability"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
)

// CreateFakeDatagram creates a fake datagram holding the text, sent from the
// address. An address that cannot be resolved leaves the sender unknown.
func CreateFakeDatagram(text, sender string) *ability.Datagram {
	return CreateFakeBinaryDatagram([]byte(text), sender)
}

// CreateFakeBinaryDatagram creates a fake datagram holding the bytes, sent
// from the address. An address that cannot be resolved leaves the sender
// unknown.
func CreateFakeBinaryDatagram(data []byte, sender string) *ability.Datagram {
	address, err := net.ResolveUDPAddr("udp", sender)
	if err != nil {
		address = nil
	}

	return impl.NewDatagram(data, address)
}
