package fixture

import (
	"errors"
	"net"
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability/communicate"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
)

// ErrConnectionRefused is the error with which the system refuses to open the
// connections of the actor created by ActorWhoCannotOpenAConnection.
var ErrConnectionRefused = errors.New("the system refuses to open the connection")

// NarratedActor creates an actor whose narration is recorded by the recorder,
// and who can communicate over UDP using the connection. With a nil
// connection, the actor can communicate over UDP but holds no connection yet.
func NarratedActor(name string, recorder *screenplayfixture.Recorder, conn *net.UDPConn) *screenplay.Actor {
	production := screenplay.NewProduction(screenplay.WithNarrator(screenplay.NewNarrator(recorder)))

	return production.SetTheStage(screenplay.CastOfStandardActors()).
		ActorNamed(name).
		WhoCan(communicate.OverUDP().UsingTheConnection(conn))
}

// ActorWhoCannotOpenAConnection creates an actor who can communicate over UDP
// but whose connections the system refuses to open, with ErrConnectionRefused.
// Whether an address can be connected to depends on the platform, so a test of
// that failure relies on this actor instead of on an address.
func ActorWhoCannotOpenAConnection(t *testing.T, name string) *screenplay.Actor {
	t.Helper()

	actor := screenplay.ActorNamed(name).WhoCan(communicate.OverUDP())
	refuse := func(*net.UDPAddr) (*net.UDPConn, error) { return nil, ErrConnectionRefused }
	require.NoError(t, impl.TheAbilityOf(actor).ToOpenConnectionsWith(refuse))

	return actor
}
