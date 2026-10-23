package connect

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/action"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
	"github.com/grandper/go-screenplay-udp/internal/verdict"
)

// ToTheUDPAddress connects the actor to the remote address: a datagram sent
// without an address goes there, and only the datagrams coming from there
// are received. It reads "connect to the UDP address 127.0.0.1:9000". It
// fails with action.ErrAddressIsEmpty when the address is empty, with an
// error wrapping ability.ErrInvalidAddress when it cannot be resolved, with
// an error wrapping ability.ErrFailedToConnect when the connection cannot be
// opened, and with ability.ErrAlreadyListening when the actor already holds
// a connection.
func ToTheUDPAddress(address string) screenplay.Performable {
	return &toTheUDPAddressAction{address: address}
}

// toTheUDPAddressAction connects the actor to a remote UDP address.
type toTheUDPAddressAction struct {
	address string
}

// PerformAs performs the action as the provided actor.
func (a *toTheUDPAddressAction) PerformAs(theActor *screenplay.Actor) error {
	if a.address == "" {
		return action.ErrAddressIsEmpty
	}

	return impl.TheAbilityOf(theActor).ToConnectTo(a.address)
}

// String describes the action.
func (a *toTheUDPAddressAction) String() string {
	return "connect to the UDP address " + verdict.Describe(a.address)
}

// toTheUDPAddressAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*toTheUDPAddressAction)(nil)
