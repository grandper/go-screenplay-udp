package listen

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/action"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
	"github.com/grandper/go-screenplay-udp/internal/verdict"
)

// OnTheUDPAddress binds the connection of the actor to the address; it reads
// "listen on the UDP address 127.0.0.1:9000". A port of 0 lets the system
// choose a free one. It fails with action.ErrAddressIsEmpty when the address
// is empty, with an error wrapping ability.ErrInvalidAddress when it cannot
// be resolved, with an error wrapping ability.ErrFailedToListen when it
// cannot be bound, and with ability.ErrAlreadyListening when the actor
// already holds a connection.
func OnTheUDPAddress(address string) screenplay.Performable {
	return &onTheUDPAddressAction{
		address:     address,
		description: "listen on the UDP address " + verdict.Describe(address),
	}
}

// OnAFreeUDPPort binds the connection of the actor to a port chosen by the
// system on the loopback interface; it reads "listen on a free UDP port". The
// question LocalUDPAddress tells which port was chosen. It fails with
// ability.ErrAlreadyListening when the actor already holds a connection.
func OnAFreeUDPPort() screenplay.Performable {
	return &onTheUDPAddressAction{
		address:     impl.FreePortAddress,
		description: "listen on a free UDP port",
	}
}

// onTheUDPAddressAction binds the connection of the actor to an address.
type onTheUDPAddressAction struct {
	address     string
	description string
}

// PerformAs performs the action as the provided actor.
func (a *onTheUDPAddressAction) PerformAs(theActor *screenplay.Actor) error {
	if a.address == "" {
		return action.ErrAddressIsEmpty
	}

	return impl.TheAbilityOf(theActor).ToListenOn(a.address)
}

// String describes the action.
func (a *onTheUDPAddressAction) String() string {
	return a.description
}

// onTheUDPAddressAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*onTheUDPAddressAction)(nil)
