package listen_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability/communicate"
	"github.com/grandper/go-screenplay-udp/action/listen"
	"github.com/grandper/go-screenplay-udp/fixture"
	internalfixture "github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/question/local"
	"github.com/grandper/go-screenplay-udp/udp"
)

const multicastGroup = "224.0.0.251:0"

func TestOnTheMulticastGroup(t *testing.T) {
	t.Run("listens on a multicast group", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())

			require.NoError(t, adam.AttemptsTo(listen.OnTheMulticastGroup(multicastGroup)))

			local, err := adam.AsksFor(local.UDPAddress())
			require.NoError(t, err)
			assert.Positive(t, local.Port)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the address is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
			require.ErrorIs(t, adam.AttemptsTo(listen.OnTheMulticastGroup("")), udp.ErrAddressIsEmpty)
		})

		t.Run("fails when the address cannot be resolved", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
			require.ErrorIs(t, adam.AttemptsTo(listen.OnTheMulticastGroup("not an address")), udp.ErrInvalidAddress)
		})

		t.Run("fails when the address is not a multicast group", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
			require.ErrorIs(t, adam.AttemptsTo(listen.OnTheMulticastGroup("127.0.0.1:0")), udp.ErrFailedToListen)
		})

		t.Run("fails when the actor already holds a connection", func(t *testing.T) {
			conn := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().UsingTheConnection(conn))
			require.ErrorIs(t, adam.AttemptsTo(listen.OnTheMulticastGroup(multicastGroup)), udp.ErrAlreadyListening)
		})

		t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(listen.OnTheMulticastGroup(multicastGroup))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		adam := internalfixture.NarratedActor("Adam", recorder, nil)

		require.NoError(t, adam.AttemptsTo(listen.OnTheMulticastGroup(multicastGroup)))

		assert.Equal(t, `Adam attempts to listen on the multicast group "224.0.0.251:0"`, recorder.Messages()[0])
		require.NoError(t, adam.Exit())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(
			t,
			`listen on the multicast group "224.0.0.251:5353"`,
			listen.OnTheMulticastGroup("224.0.0.251:5353").String(),
		)
	})
}
