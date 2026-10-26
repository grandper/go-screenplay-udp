package stop

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-udp/internal/ability"
)

// ListeningOverUDP closes the connection of the actor and forgets the address
// it was bound to, so that the actor can listen again; the datagrams received
// so far are kept. It reads "stop listening over UDP". It fails with
// udp.ErrUDPConnectionNotSet when the actor holds no connection.
func ListeningOverUDP() screenplay.Performable {
	return &listeningOverUDPAction{}
}

// listeningOverUDPAction closes the connection of the actor.
type listeningOverUDPAction struct{}

// PerformAs performs the action as the provided actor.
func (a *listeningOverUDPAction) PerformAs(theActor *screenplay.Actor) error {
	return impl.TheAbilityOf(theActor).ToClose()
}

// String describes the action.
func (a *listeningOverUDPAction) String() string {
	return "stop listening over UDP"
}

// listeningOverUDPAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*listeningOverUDPAction)(nil)
