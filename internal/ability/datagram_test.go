package ability_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	impl "github.com/grandper/go-screenplay-udp/internal/ability"
)

func TestDatagram(t *testing.T) {
	sender := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000}
	data := []byte("Hello")

	t.Run("holds a copy of the payload", func(t *testing.T) {
		datagram := impl.NewDatagram(data, sender)
		data[0] = 'J'

		assert.Equal(t, []byte("Hello"), datagram.Data())
		assert.Equal(t, "Hello", datagram.Text())
		data[0] = 'H'
	})

	t.Run("returns a copy of the payload", func(t *testing.T) {
		datagram := impl.NewDatagram(data, sender)
		datagram.Data()[0] = 'J'

		assert.Equal(t, "Hello", datagram.Text())
	})

	t.Run("holds a copy of the sender", func(t *testing.T) {
		datagram := impl.NewDatagram(data, sender)
		sender.Port = 9001
		defer func() { sender.Port = 9000 }()

		require.NotNil(t, datagram.Sender())
		assert.Equal(t, 9000, datagram.Sender().Port)
		datagram.Sender().Port = 9002
		assert.Equal(t, 9000, datagram.Sender().Port)
	})

	t.Run("has no sender when it is unknown", func(t *testing.T) {
		datagram := impl.NewDatagram(data, nil)

		assert.Nil(t, datagram.Sender())
		assert.False(t, datagram.IsFrom(sender))
	})

	t.Run("tells if it was sent from an address", func(t *testing.T) {
		datagram := impl.NewDatagram(data, sender)

		assert.True(t, datagram.IsFrom(&net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000}))
		assert.False(t, datagram.IsFrom(&net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9001}))
		assert.False(t, datagram.IsFrom(nil))
	})
}
