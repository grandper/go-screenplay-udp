package fixture

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	// maxDatagramSize is the largest payload a UDP datagram can carry.
	maxDatagramSize = 65535

	// freePortAddress is a port chosen by the system on the loopback interface.
	freePortAddress = "127.0.0.1:0"
)

// ListenOnAFreeUDPPort binds a connection to a port chosen by the system on
// the loopback interface, closed when the test ends, and returns it.
func ListenOnAFreeUDPPort(t *testing.T) *net.UDPConn {
	t.Helper()

	local, err := net.ResolveUDPAddr("udp", freePortAddress)
	require.NoError(t, err)
	conn, err := net.ListenUDP("udp", local)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	return conn
}

// RunAnEchoUDPServer runs a server on a free port of the loopback interface
// that sends every datagram it receives back to its sender, stopped when the
// test ends, and returns its address.
func RunAnEchoUDPServer(t *testing.T) *net.UDPAddr {
	t.Helper()

	return RunAUDPServer(t, func(data []byte) []byte { return data })
}

// RunAUDPServer runs a server on a free port of the loopback interface that
// replies to every datagram it receives with the reply the function returns,
// or does not reply when the function returns nil. It is stopped when the
// test ends, and its address is returned.
func RunAUDPServer(t *testing.T, reply func(data []byte) []byte) *net.UDPAddr {
	t.Helper()

	conn := ListenOnAFreeUDPPort(t)
	address, ok := conn.LocalAddr().(*net.UDPAddr)
	require.True(t, ok)

	go serve(conn, reply)

	return address
}

// serve answers the datagrams received on the connection until it is closed,
// or until a read fails: a read that keeps failing, for example once a
// deadline is set on the connection, must not spin.
func serve(conn *net.UDPConn, reply func(data []byte) []byte) {
	buffer := make([]byte, maxDatagramSize)
	for {
		read, sender, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return
		}
		if response := reply(buffer[:read]); response != nil {
			_, _ = conn.WriteToUDP(response, sender)
		}
	}
}
