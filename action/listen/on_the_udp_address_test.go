package listen_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability/communicate"
	"github.com/grandper/go-screenplay-udp/action/listen"
	"github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/question/local"
	"github.com/grandper/go-screenplay-udp/udp"
)

func TestOnTheUDPAddress(t *testing.T) {
	t.Run("listens on a UDP address", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())

			require.NoError(t, adam.AttemptsTo(listen.OnTheUDPAddress("127.0.0.1:0")))

			local, err := adam.AsksFor(local.UDPAddress())
			require.NoError(t, err)
			assert.Equal(t, "127.0.0.1", local.IP.String())
			assert.Positive(t, local.Port)
			require.NoError(t, adam.Exit())
		})

		t.Run("on a free port", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())

			require.NoError(t, adam.AttemptsTo(listen.OnAFreeUDPPort()))

			local, err := adam.AsksFor(local.UDPAddress())
			require.NoError(t, err)
			assert.Equal(t, "127.0.0.1", local.IP.String())
			assert.Positive(t, local.Port)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the address is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
			require.ErrorIs(t, adam.AttemptsTo(listen.OnTheUDPAddress("")), udp.ErrAddressIsEmpty)
		})

		t.Run("fails when the address cannot be resolved", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
			require.ErrorIs(t, adam.AttemptsTo(listen.OnTheUDPAddress("not an address")), udp.ErrInvalidAddress)
		})

		t.Run("fails when the address cannot be bound", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
			require.ErrorIs(t, adam.AttemptsTo(listen.OnTheUDPAddress("192.0.2.1:0")), udp.ErrFailedToListen)
		})

		t.Run("fails when the actor already holds a connection", func(t *testing.T) {
			conn := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().UsingTheConnection(conn))
			require.ErrorIs(t, adam.AttemptsTo(listen.OnTheUDPAddress("127.0.0.1:0")), udp.ErrAlreadyListening)
			require.ErrorIs(t, adam.AttemptsTo(listen.OnAFreeUDPPort()), udp.ErrAlreadyListening)
		})

		t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(listen.OnTheUDPAddress("127.0.0.1:0"))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
			err = screenplay.ActorNamed("Adam").AttemptsTo(listen.OnAFreeUDPPort())
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		adam := fixture.NarratedActor("Adam", recorder, nil)

		require.NoError(t, adam.AttemptsTo(listen.OnTheUDPAddress("127.0.0.1:0")))

		assert.Equal(t, `Adam attempts to listen on the UDP address "127.0.0.1:0"`, recorder.Messages()[0])
		require.NoError(t, adam.Exit())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `listen on the UDP address "127.0.0.1:9000"`, listen.OnTheUDPAddress("127.0.0.1:9000").String())
		assert.Equal(t, "listen on a free UDP port", listen.OnAFreeUDPPort().String())
	})
}
