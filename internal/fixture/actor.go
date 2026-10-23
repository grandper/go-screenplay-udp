package fixture

import (
	"net"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"

	"github.com/grandper/go-screenplay-udp/ability"
)

// NarratedActor creates an actor whose narration is recorded by the recorder,
// and who can communicate over UDP using the connection. With a nil
// connection, the actor can communicate over UDP but holds no connection yet.
func NarratedActor(name string, recorder *screenplayfixture.Recorder, conn *net.UDPConn) *screenplay.Actor {
	production := screenplay.NewProduction(screenplay.WithNarrator(screenplay.NewNarrator(recorder)))

	return production.SetTheStage(screenplay.CastOfStandardActors()).
		ActorNamed(name).
		WhoCan(ability.CommunicateOverUDP().UsingTheConnection(conn))
}
