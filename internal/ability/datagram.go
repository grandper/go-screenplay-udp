package ability

import (
	"net"
	"slices"
)

// Datagram is a datagram received by the actor: its payload and the address
// it was sent from. It is read through Data, Text and Sender, so a test cannot
// change it.
type Datagram struct {
	data   []byte
	sender *net.UDPAddr
}

// NewDatagram creates a datagram holding a copy of the payload, sent from the
// address.
func NewDatagram(data []byte, sender *net.UDPAddr) *Datagram {
	return &Datagram{data: slices.Clone(data), sender: cloneAddress(sender)}
}

// Data returns a copy of the payload of the datagram.
func (d *Datagram) Data() []byte {
	return slices.Clone(d.data)
}

// Text returns the payload of the datagram as a text.
func (d *Datagram) Text() string {
	return string(d.data)
}

// Sender returns a copy of the address the datagram was sent from, nil when
// it is unknown.
func (d *Datagram) Sender() *net.UDPAddr {
	return cloneAddress(d.sender)
}

// IsFrom tells if the datagram was sent from the address.
func (d *Datagram) IsFrom(address *net.UDPAddr) bool {
	if d.sender == nil || address == nil {
		return false
	}

	return d.sender.String() == address.String()
}

// cloneAddress copies the address, so that a caller cannot change the one the
// ability holds.
func cloneAddress(address *net.UDPAddr) *net.UDPAddr {
	if address == nil {
		return nil
	}

	return &net.UDPAddr{
		IP:   slices.Clone(address.IP),
		Port: address.Port,
		Zone: address.Zone,
	}
}
