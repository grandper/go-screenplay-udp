package listen

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-udp/internal/ability"
	"github.com/grandper/go-screenplay-udp/internal/verdict"
	"github.com/grandper/go-screenplay-udp/udp"
)

// OnTheMulticastGroup binds the connection of the actor to the address of a
// multicast group, on the multicast interface of the system, so that the
// actor receives the datagrams sent to the group; it reads "listen on the
// multicast group 224.0.0.251:5353". It fails with udp.ErrAddressIsEmpty
// when the address is empty, with an error wrapping
// udp.ErrInvalidAddress when it cannot be resolved, with an error
// wrapping udp.ErrFailedToListen when the group cannot be joined, and
// with udp.ErrAlreadyListening when the actor already holds a connection.
func OnTheMulticastGroup(address string) screenplay.Performable {
	return &onTheMulticastGroupAction{address: address}
}

// onTheMulticastGroupAction binds the connection of the actor to a multicast
// group.
type onTheMulticastGroupAction struct {
	address string
}

// PerformAs performs the action as the provided actor.
func (a *onTheMulticastGroupAction) PerformAs(theActor *screenplay.Actor) error {
	if a.address == "" {
		return udp.ErrAddressIsEmpty
	}

	return impl.TheAbilityOf(theActor).ToListenOnTheMulticastGroup(a.address)
}

// String describes the action.
func (a *onTheMulticastGroupAction) String() string {
	return "listen on the multicast group " + verdict.Describe(a.address)
}

// onTheMulticastGroupAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*onTheMulticastGroupAction)(nil)
