package udp_test

import (
	"net"
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability/communicate"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
	"github.com/grandper/go-screenplay-udp/fixture"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
	"github.com/grandper/go-screenplay-udp/udp"
)

func TestDatagram(t *testing.T) {
	sender := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000}

	t.Run("is the datagram of the ability", func(t *testing.T) {
		datagram := impl.NewDatagram([]byte("Hello"), sender)

		assert.IsType(t, (*udp.Datagram)(nil), datagram)
		assert.Equal(t, []byte("Hello"), datagram.Data())
		assert.Equal(t, "Hello", datagram.Text())
		require.NotNil(t, datagram.Sender())
		assert.Equal(t, "127.0.0.1:9000", datagram.Sender().String())
		assert.True(t, datagram.IsFrom(sender))
	})

	t.Run("is what the actor receives", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
		require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping").To(echo.String()), receive.TheNextDatagram()))

		datagrams, err := adam.AsksFor(udp.Datagrams())

		require.NoError(t, err)
		require.Len(t, datagrams, 1)
		assert.IsType(t, (*udp.Datagram)(nil), datagrams[0])
		assert.Equal(t, "ping", datagrams[0].Text())
		assert.Equal(t, echo.String(), datagrams[0].Sender().String())
		require.NoError(t, adam.Exit())
	})

	t.Run("does not change once received", func(t *testing.T) {
		data := []byte("Hello")
		datagram := impl.NewDatagram(data, sender)

		data[0] = 'J'
		datagram.Data()[0] = 'J'
		datagram.Sender().Port = 9001

		assert.Equal(t, "Hello", datagram.Text())
		assert.Equal(t, []byte("Hello"), datagram.Data())
		assert.Equal(t, 9000, datagram.Sender().Port)
	})
}
