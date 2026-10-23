package question

import (
	"net"

	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-udp/internal/ability"
)

// LocalUDPAddress asks for the address the connection of the actor is bound
// to, which tells the port the system chose when the actor listens on a free
// port.
func LocalUDPAddress() screenplay.Question[*net.UDPAddr] {
	return &localUDPAddressQuestion{}
}

// localUDPAddressQuestion asks for the local address of the connection.
type localUDPAddressQuestion struct{}

// AnsweredBy returns the address the connection of the actor is bound to.
func (q *localUDPAddressQuestion) AnsweredBy(theActor *screenplay.Actor) (*net.UDPAddr, error) {
	return impl.TheAbilityOf(theActor).ToTellTheLocalAddress()
}

// String describes the question.
func (q *localUDPAddressQuestion) String() string {
	return "local UDP address"
}

// localUDPAddressQuestion implements the screenplay.Question interface.
var _ screenplay.Question[*net.UDPAddr] = (*localUDPAddressQuestion)(nil)
