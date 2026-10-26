package udp_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/question/number"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability/communicate"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
	"github.com/grandper/go-screenplay-udp/fixture"
	"github.com/grandper/go-screenplay-udp/udp"
)

func TestDatagrams(t *testing.T) {
	t.Run("answers with the datagrams received by the actor", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		other := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
		require.NoError(t, adam.AttemptsTo(
			send.TheDatagram("1").To(echo.String()),
			receive.TheNextDatagram(),
			send.TheDatagram("2").To(other.String()),
			receive.TheNextDatagram(),
			send.TheDatagram("3").To(echo.String()),
			receive.TheNextDatagram(),
		))

		t.Run("all of them, in the order they were received", func(t *testing.T) {
			datagrams, err := adam.AsksFor(udp.Datagrams())

			require.NoError(t, err)
			require.Len(t, datagrams, 3)
			assert.Equal(t, "1", datagrams[0].Text())
			assert.Equal(t, "2", datagrams[1].Text())
			assert.Equal(t, "3", datagrams[2].Text())
		})

		t.Run("those from an address", func(t *testing.T) {
			datagrams, err := adam.AsksFor(udp.DatagramsFrom(echo.String()))

			require.NoError(t, err)
			require.Len(t, datagrams, 2)
			assert.Equal(t, "1", datagrams[0].Text())
			assert.Equal(t, "3", datagrams[1].Text())
		})

		t.Run("those from an address written with a host name", func(t *testing.T) {
			datagrams, err := adam.AsksFor(udp.DatagramsFrom("localhost:" + portOf(other)))

			require.NoError(t, err)
			require.Len(t, datagrams, 1)
			assert.Equal(t, "2", datagrams[0].Text())
		})

		t.Run("none from an address nothing was received from", func(t *testing.T) {
			datagrams, err := adam.AsksFor(udp.DatagramsFrom("127.0.0.1:9"))

			require.NoError(t, err)
			assert.Empty(t, datagrams)
		})

		t.Run("fails when the address cannot be resolved", func(t *testing.T) {
			_, err := adam.AsksFor(udp.DatagramsFrom("not an address"))

			require.ErrorIs(t, err, udp.ErrInvalidAddress)
		})

		t.Run("can be composed with the core questions", func(t *testing.T) {
			require.NoError(t, adam.Should(see.The(number.Of(udp.Datagrams())).Should(be.EqualTo(3))))
			lastDatagram, err := adam.AsksFor(last.Of(udp.DatagramsFrom(echo.String())))
			require.NoError(t, err)
			assert.Equal(t, "3", lastDatagram.Text())
		})

		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
		_, err := screenplay.ActorNamed("Adam").AsksFor(udp.Datagrams())
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		_, err = screenplay.ActorNamed("Adam").AsksFor(udp.DatagramsFrom("127.0.0.1:9000"))
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "UDP datagrams", udp.Datagrams().String())
		assert.Equal(t, `UDP datagrams from "127.0.0.1:9000"`, udp.DatagramsFrom("127.0.0.1:9000").String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[[]*udp.Datagram])(nil), udp.Datagrams())
		assert.Implements(
			t,
			(*screenplay.Question[[]*udp.Datagram])(nil),
			udp.DatagramsFrom("127.0.0.1:9000"),
		)
	})
}

// portOf returns the port of the address as a text.
func portOf(address *net.UDPAddr) string {
	return strconv.Itoa(address.Port)
}
