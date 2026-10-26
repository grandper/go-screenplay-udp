package local

import (
	"net"

	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-udp/internal/ability"
)

// UDPAddress asks for the address the connection of the actor is bound
// to, which tells the port the system chose when the actor listens on a free
// port.
func UDPAddress() screenplay.Question[*net.UDPAddr] {
	return &udpAddressQuestion{}
}

// udpAddressQuestion asks for the local address of the connection.
type udpAddressQuestion struct{}

// AnsweredBy returns the address the connection of the actor is bound to.
func (q *udpAddressQuestion) AnsweredBy(theActor *screenplay.Actor) (*net.UDPAddr, error) {
	return impl.TheAbilityOf(theActor).ToTellTheLocalAddress()
}

// String describes the question.
func (q *udpAddressQuestion) String() string {
	return "local UDP address"
}

// udpAddressQuestion implements the screenplay.Question interface.
var _ screenplay.Question[*net.UDPAddr] = (*udpAddressQuestion)(nil)
