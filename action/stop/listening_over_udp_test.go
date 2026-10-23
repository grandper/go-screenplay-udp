package stop_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/question/number"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/action/listen"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
	"github.com/grandper/go-screenplay-udp/action/stop"
	"github.com/grandper/go-screenplay-udp/fixture"
	internalfixture "github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/question"
)

func TestListeningOverUDP(t *testing.T) {
	t.Run("stops listening over UDP", func(t *testing.T) {
		t.Run("and lets the actor listen again, keeping the datagrams", func(t *testing.T) {
			echo := fixture.RunAnEchoUDPServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping").To(echo.String()), receive.TheNextDatagram()))

			require.NoError(t, adam.AttemptsTo(stop.ListeningOverUDP()))

			err := adam.AttemptsTo(send.TheDatagram("ping").To(echo.String()))
			require.ErrorIs(t, err, ability.ErrUDPConnectionNotSet)
			require.NoError(t, adam.AttemptsTo(listen.OnAFreeUDPPort()))
			require.NoError(t, adam.Should(see.The(number.Of(question.UDPDatagrams())).Should(be.EqualTo(1))))
			require.NoError(t, adam.Exit())
		})

		t.Run("closes a connection the actor was given", func(t *testing.T) {
			conn := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().UsingTheConnection(conn))

			require.NoError(t, adam.AttemptsTo(stop.ListeningOverUDP()))

			_, err := conn.WriteToUDP([]byte("ping"), fixture.RunAnEchoUDPServer(t))
			require.Error(t, err)
		})

		t.Run("fails when the actor holds no connection", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP())
			require.ErrorIs(t, adam.AttemptsTo(stop.ListeningOverUDP()), ability.ErrUDPConnectionNotSet)
		})

		t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(stop.ListeningOverUDP())
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		adam := internalfixture.NarratedActor("Adam", recorder, fixture.ListenOnAFreeUDPPort(t))

		require.NoError(t, adam.AttemptsTo(stop.ListeningOverUDP()))

		assert.Equal(t, "Adam attempts to stop listening over UDP", recorder.Messages()[0])
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "stop listening over UDP", stop.ListeningOverUDP().String())
	})
}
