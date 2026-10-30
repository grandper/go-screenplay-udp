package fixture_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/internal/fixture"
)

const timeout = 2 * time.Second

func TestListenOnAFreeUDPPort(t *testing.T) {
	first := fixture.ListenOnAFreeUDPPort(t)
	second := fixture.ListenOnAFreeUDPPort(t)

	assert.NotEqual(t, first.LocalAddr().String(), second.LocalAddr().String())
	local, ok := first.LocalAddr().(*net.UDPAddr)
	require.True(t, ok)
	assert.Equal(t, "127.0.0.1", local.IP.String())
	assert.Positive(t, local.Port)
}

func TestRunAnEchoUDPServer(t *testing.T) {
	echo := fixture.RunAnEchoUDPServer(t)
	client := fixture.ListenOnAFreeUDPPort(t)

	_, err := client.WriteToUDP([]byte("ping"), echo)
	require.NoError(t, err)

	assert.Equal(t, "ping", readFrom(t, client))
}

func TestRunAUDPServer(t *testing.T) {
	t.Run("replies with what the function returns", func(t *testing.T) {
		server := fixture.RunAUDPServer(t, func(data []byte) []byte { return append([]byte("re: "), data...) })
		client := fixture.ListenOnAFreeUDPPort(t)

		_, err := client.WriteToUDP([]byte("ping"), server)
		require.NoError(t, err)

		assert.Equal(t, "re: ping", readFrom(t, client))
	})

	t.Run("does not reply when the function returns nil", func(t *testing.T) {
		server := fixture.RunAUDPServer(t, func([]byte) []byte { return nil })
		client := fixture.ListenOnAFreeUDPPort(t)

		_, err := client.WriteToUDP([]byte("ping"), server)
		require.NoError(t, err)

		require.NoError(t, client.SetReadDeadline(time.Now().Add(100*time.Millisecond)))
		_, _, err = client.ReadFromUDP(make([]byte, 1024))
		require.Error(t, err)
	})
}

// readFrom reads the next datagram of the connection as a text.
func readFrom(t *testing.T, conn *net.UDPConn) string {
	t.Helper()

	require.NoError(t, conn.SetReadDeadline(time.Now().Add(timeout)))
	buffer := make([]byte, 1024)
	read, _, err := conn.ReadFromUDP(buffer)
	require.NoError(t, err)

	return string(buffer[:read])
}
