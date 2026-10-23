package fixture_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/action/listen"
	"github.com/grandper/go-screenplay-udp/action/send"
	udpfixture "github.com/grandper/go-screenplay-udp/fixture"
	"github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/question"
)

func TestNarratedActor(t *testing.T) {
	t.Run("creates an actor with the name", func(t *testing.T) {
		adam := fixture.NarratedActor("Adam", screenplayfixture.NewRecorder(), nil)

		assert.Equal(t, "Adam", adam.Name())
	})

	t.Run("records the narration of the actor", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		peer := udpfixture.ListenOnAFreeUDPPort(t)
		adam := fixture.NarratedActor("Adam", recorder, udpfixture.ListenOnAFreeUDPPort(t))

		require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping").To(peer.LocalAddr().String())))

		require.NotEmpty(t, recorder.Messages())
		assert.Equal(
			t,
			`Adam attempts to send the datagram "ping" to "`+peer.LocalAddr().String()+`"`,
			recorder.Messages()[0],
		)
	})

	t.Run("gives the actor the ability to communicate over UDP using the connection", func(t *testing.T) {
		conn := udpfixture.ListenOnAFreeUDPPort(t)
		adam := fixture.NarratedActor("Adam", screenplayfixture.NewRecorder(), conn)

		assert.True(t, adam.HasAbilityTo(ability.CommunicateOverUDP()))
		local, err := adam.AsksFor(question.LocalUDPAddress())
		require.NoError(t, err)
		assert.Equal(t, conn.LocalAddr().String(), local.String())
	})

	t.Run("gives the actor the ability without a connection when it is nil", func(t *testing.T) {
		adam := fixture.NarratedActor("Adam", screenplayfixture.NewRecorder(), nil)

		assert.True(t, adam.HasAbilityTo(ability.CommunicateOverUDP()))
		_, err := adam.AsksFor(question.LocalUDPAddress())
		require.ErrorIs(t, err, ability.ErrUDPConnectionNotSet)

		require.NoError(t, adam.AttemptsTo(listen.OnAFreeUDPPort()))
		_, err = adam.AsksFor(question.LocalUDPAddress())
		require.NoError(t, err)
		require.NoError(t, adam.Exit())
	})

	t.Run("creates independent actors", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		adam := fixture.NarratedActor("Adam", recorder, udpfixture.ListenOnAFreeUDPPort(t))
		eve := fixture.NarratedActor("Eve", recorder, udpfixture.ListenOnAFreeUDPPort(t))

		adamAddress, err := adam.AsksFor(question.LocalUDPAddress())
		require.NoError(t, err)
		eveAddress, err := eve.AsksFor(question.LocalUDPAddress())
		require.NoError(t, err)

		assert.NotEqual(t, adamAddress.String(), eveAddress.String())
		require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping").To(eveAddress.String())))
		assert.Contains(t, recorder.Messages(), "Adam asks for the local UDP address")
		assert.Contains(t, recorder.Messages(), "Eve asks for the local UDP address")
		assert.Contains(
			t,
			recorder.Messages(),
			`Adam attempts to send the datagram "ping" to "`+eveAddress.String()+`"`,
		)
	})

	t.Run("implements the actor of go-screenplay", func(t *testing.T) {
		assert.IsType(t, (*screenplay.Actor)(nil), fixture.NarratedActor("Adam", screenplayfixture.NewRecorder(), nil))
	})
}
