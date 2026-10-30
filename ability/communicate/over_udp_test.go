package communicate_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability/communicate"
	"github.com/grandper/go-screenplay-udp/action/send"
	"github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/question/local"
	"github.com/grandper/go-screenplay-udp/udp"
)

func TestCommunicateOverUDP(t *testing.T) {
	t.Run("gives an actor the ability to communicate over UDP using a connection", func(t *testing.T) {
		conn := fixture.ListenOnAFreeUDPPort(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().UsingTheConnection(conn))

		local, err := adam.AsksFor(local.UDPAddress())

		require.NoError(t, err)
		assert.Equal(t, conn.LocalAddr().String(), local.String())
		require.NoError(t, adam.Exit())
		_, err = conn.WriteToUDP([]byte("ping"), local)
		require.Error(t, err)
	})

	t.Run("gives an actor the ability to communicate over UDP listening on an address", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOn("127.0.0.1:0"))

		local, err := adam.AsksFor(local.UDPAddress())

		require.NoError(t, err)
		assert.Equal(t, "127.0.0.1", local.IP.String())
		assert.Positive(t, local.Port)
		require.NoError(t, adam.Exit())
	})

	t.Run("gives an actor the ability to communicate over UDP listening on a free port", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())

		require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping").To(echo.String())))
		require.NoError(t, adam.Exit())
	})

	t.Run("fails every action until the actor is given a connection", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())

		err := adam.AttemptsTo(send.TheDatagram("ping").To("127.0.0.1:9000"))

		require.ErrorIs(t, err, udp.ErrUDPConnectionNotSet)
	})
}
