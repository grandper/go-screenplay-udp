package connect_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability/communicate"
	"github.com/grandper/go-screenplay-udp/action/connect"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
	"github.com/grandper/go-screenplay-udp/fixture"
	internalfixture "github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/resolution/been"
	"github.com/grandper/go-screenplay-udp/udp"
)

func TestToTheUDPAddress(t *testing.T) {
	t.Run("connects to a UDP address", func(t *testing.T) {
		t.Run("successfully, and sends to it without an address", func(t *testing.T) {
			echo := fixture.RunAnEchoUDPServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())

			require.NoError(t, adam.AttemptsTo(connect.ToTheUDPAddress(echo.String())))

			require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping"), receive.TheNextDatagram()))
			require.NoError(
				t,
				adam.Should(see.The(udp.Datagrams()).ShouldHave(been.ReceivedFrom(echo.String()))),
			)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the address is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
			require.ErrorIs(t, adam.AttemptsTo(connect.ToTheUDPAddress("")), udp.ErrAddressIsEmpty)
		})

		t.Run("fails when the address cannot be resolved", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
			require.ErrorIs(t, adam.AttemptsTo(connect.ToTheUDPAddress("not an address")), udp.ErrInvalidAddress)
		})

		t.Run("fails when the connection cannot be opened", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
			require.ErrorIs(t, adam.AttemptsTo(connect.ToTheUDPAddress("127.0.0.1:0")), udp.ErrFailedToConnect)
		})

		t.Run("fails when the actor already holds a connection", func(t *testing.T) {
			conn := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().UsingTheConnection(conn))
			require.ErrorIs(t, adam.AttemptsTo(connect.ToTheUDPAddress("127.0.0.1:9000")), udp.ErrAlreadyListening)
		})

		t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(connect.ToTheUDPAddress("127.0.0.1:9000"))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		echo := fixture.RunAnEchoUDPServer(t)
		adam := internalfixture.NarratedActor("Adam", recorder, nil)

		require.NoError(t, adam.AttemptsTo(connect.ToTheUDPAddress(echo.String())))

		assert.Equal(t, `Adam attempts to connect to the UDP address "`+echo.String()+`"`, recorder.Messages()[0])
		require.NoError(t, adam.Exit())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(
			t,
			`connect to the UDP address "127.0.0.1:9000"`,
			connect.ToTheUDPAddress("127.0.0.1:9000").String(),
		)
	})
}
