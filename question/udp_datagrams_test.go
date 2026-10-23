package question_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/question/number"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
	"github.com/grandper/go-screenplay-udp/fixture"
	"github.com/grandper/go-screenplay-udp/question"
)

func TestUDPDatagrams(t *testing.T) {
	t.Run("answers with the datagrams received by the actor", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		other := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
		require.NoError(t, adam.AttemptsTo(
			send.TheDatagram("1").To(echo.String()),
			receive.TheNextDatagram(),
			send.TheDatagram("2").To(other.String()),
			receive.TheNextDatagram(),
			send.TheDatagram("3").To(echo.String()),
			receive.TheNextDatagram(),
		))

		t.Run("all of them, in the order they were received", func(t *testing.T) {
			datagrams, err := adam.AsksFor(question.UDPDatagrams())

			require.NoError(t, err)
			require.Len(t, datagrams, 3)
			assert.Equal(t, "1", datagrams[0].Text())
			assert.Equal(t, "2", datagrams[1].Text())
			assert.Equal(t, "3", datagrams[2].Text())
		})

		t.Run("those from an address", func(t *testing.T) {
			datagrams, err := adam.AsksFor(question.UDPDatagramsFrom(echo.String()))

			require.NoError(t, err)
			require.Len(t, datagrams, 2)
			assert.Equal(t, "1", datagrams[0].Text())
			assert.Equal(t, "3", datagrams[1].Text())
		})

		t.Run("those from an address written with a host name", func(t *testing.T) {
			datagrams, err := adam.AsksFor(question.UDPDatagramsFrom("localhost:" + portOf(other)))

			require.NoError(t, err)
			require.Len(t, datagrams, 1)
			assert.Equal(t, "2", datagrams[0].Text())
		})

		t.Run("none from an address nothing was received from", func(t *testing.T) {
			datagrams, err := adam.AsksFor(question.UDPDatagramsFrom("127.0.0.1:9"))

			require.NoError(t, err)
			assert.Empty(t, datagrams)
		})

		t.Run("fails when the address cannot be resolved", func(t *testing.T) {
			_, err := adam.AsksFor(question.UDPDatagramsFrom("not an address"))

			require.ErrorIs(t, err, ability.ErrInvalidAddress)
		})

		t.Run("can be composed with the core questions", func(t *testing.T) {
			require.NoError(t, adam.Should(see.The(number.Of(question.UDPDatagrams())).Should(be.EqualTo(3))))
			lastDatagram, err := adam.AsksFor(last.Of(question.UDPDatagramsFrom(echo.String())))
			require.NoError(t, err)
			assert.Equal(t, "3", lastDatagram.Text())
		})

		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
		_, err := screenplay.ActorNamed("Adam").AsksFor(question.UDPDatagrams())
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		_, err = screenplay.ActorNamed("Adam").AsksFor(question.UDPDatagramsFrom("127.0.0.1:9000"))
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "UDP datagrams", question.UDPDatagrams().String())
		assert.Equal(t, `UDP datagrams from "127.0.0.1:9000"`, question.UDPDatagramsFrom("127.0.0.1:9000").String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[[]*ability.Datagram])(nil), question.UDPDatagrams())
		assert.Implements(
			t,
			(*screenplay.Question[[]*ability.Datagram])(nil),
			question.UDPDatagramsFrom("127.0.0.1:9000"),
		)
	})
}
