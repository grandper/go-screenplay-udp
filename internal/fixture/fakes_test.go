package fixture_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/question/fixture"
)

func TestCreateFakeDatagram(t *testing.T) {
	t.Run("creates a datagram holding a text, sent from an address", func(t *testing.T) {
		datagram := fixture.CreateFakeDatagram("pong", "127.0.0.1:9000")

		assert.Equal(t, "pong", datagram.Text())
		assert.Equal(t, []byte("pong"), datagram.Data())
		require.NotNil(t, datagram.Sender())
		assert.Equal(t, "127.0.0.1:9000", datagram.Sender().String())
	})

	t.Run("creates a datagram holding bytes", func(t *testing.T) {
		datagram := fixture.CreateFakeBinaryDatagram([]byte{0x00, 0x01}, "127.0.0.1:9000")

		assert.Equal(t, []byte{0x00, 0x01}, datagram.Data())
	})

	t.Run("leaves the sender unknown when the address cannot be resolved", func(t *testing.T) {
		datagram := fixture.CreateFakeDatagram("pong", "not an address")

		assert.Nil(t, datagram.Sender())
	})
}
