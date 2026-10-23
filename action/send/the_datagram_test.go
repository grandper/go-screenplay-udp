package send_test

import (
	"net"
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/action"
	"github.com/grandper/go-screenplay-udp/action/connect"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
	"github.com/grandper/go-screenplay-udp/fixture"
	internalfixture "github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/question"
)

const (
	address = "127.0.0.1:9000"
	text    = "Hello, UDP!"
	timeout = 2 * time.Second
)

func TestTheDatagram(t *testing.T) {
	binary := []byte{0x00, 0x01, 0x02}

	t.Run("sends a datagram", func(t *testing.T) {
		t.Run("with a text", func(t *testing.T) {
			peer := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())

			require.NoError(t, adam.AttemptsTo(send.TheDatagram(text).To(peer.LocalAddr().String())))

			data, sender := readFrom(t, peer)
			assert.Equal(t, []byte(text), data)
			local, err := adam.AsksFor(question.LocalUDPAddress())
			require.NoError(t, err)
			assert.Equal(t, local.String(), sender.String())
			require.NoError(t, adam.Exit())
		})

		t.Run("with binary content", func(t *testing.T) {
			peer := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())

			require.NoError(t, adam.AttemptsTo(send.TheBinaryDatagram(binary).To(peer.LocalAddr().String())))

			data, _ := readFrom(t, peer)
			assert.Equal(t, binary, data)
			require.NoError(t, adam.Exit())
		})

		t.Run("with no content", func(t *testing.T) {
			peer := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())

			require.NoError(t, adam.AttemptsTo(send.AnEmptyDatagram().To(peer.LocalAddr().String())))

			data, _ := readFrom(t, peer)
			assert.Empty(t, data)
			require.NoError(t, adam.Exit())
		})

		t.Run("to the sender of the last datagram", func(t *testing.T) {
			peer := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			local, err := adam.AsksFor(question.LocalUDPAddress())
			require.NoError(t, err)
			_, err = peer.WriteToUDP([]byte("ping"), local)
			require.NoError(t, err)
			require.NoError(t, adam.AttemptsTo(receive.TheNextDatagram()))

			require.NoError(t, adam.AttemptsTo(send.TheDatagram("pong").ToTheSenderOfTheLastDatagram()))

			data, _ := readFrom(t, peer)
			assert.Equal(t, []byte("pong"), data)
			require.NoError(t, adam.Exit())
		})

		t.Run("to the address the actor connected to", func(t *testing.T) {
			peer := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP())
			require.NoError(t, adam.AttemptsTo(connect.ToTheUDPAddress(peer.LocalAddr().String())))

			require.NoError(t, adam.AttemptsTo(send.TheDatagram(text)))

			data, _ := readFrom(t, peer)
			assert.Equal(t, []byte(text), data)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the address is empty", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			require.ErrorIs(t, adam.AttemptsTo(send.TheDatagram(text).To("")), action.ErrAddressIsEmpty)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when no address is given to an actor that is not connected", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			require.ErrorIs(t, adam.AttemptsTo(send.TheDatagram(text)), ability.ErrDestinationNotSet)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the address cannot be resolved", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			require.ErrorIs(t, adam.AttemptsTo(send.TheDatagram(text).To("not an address")), ability.ErrInvalidAddress)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the datagram is too large", func(t *testing.T) {
			peer := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			err := adam.AttemptsTo(send.TheBinaryDatagram(make([]byte, 70000)).To(peer.LocalAddr().String()))
			require.ErrorIs(t, err, ability.ErrFailedToSend)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails to reply when no datagram has been received", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			err := adam.AttemptsTo(send.TheDatagram("pong").ToTheSenderOfTheLastDatagram())
			require.ErrorIs(t, err, ability.ErrNoDatagramReceived)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP())
			require.ErrorIs(t, adam.AttemptsTo(send.TheDatagram(text).To(address)), ability.ErrUDPConnectionNotSet)
		})

		t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(send.TheDatagram(text).To(address))
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
			err = screenplay.ActorNamed("Adam").AttemptsTo(send.TheDatagram(text).ToTheSenderOfTheLastDatagram())
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		peer := fixture.ListenOnAFreeUDPPort(t)
		adam := internalfixture.NarratedActor("Adam", recorder, fixture.ListenOnAFreeUDPPort(t))

		require.NoError(t, adam.AttemptsTo(send.TheDatagram(text).To(peer.LocalAddr().String())))

		assert.Equal(
			t,
			`Adam attempts to send the datagram "Hello, UDP!" to "`+peer.LocalAddr().String()+`"`,
			recorder.Messages()[0],
		)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `send the datagram "Hello, UDP!"`, send.TheDatagram(text).String())
		assert.Equal(
			t,
			`send the datagram "Hello, UDP!" to "127.0.0.1:9000"`,
			send.TheDatagram(text).To(address).String(),
		)
		assert.Equal(
			t,
			`send the binary datagram "\x00\x01\x02" to "127.0.0.1:9000"`,
			send.TheBinaryDatagram(binary).To(address).String(),
		)
		assert.Equal(t, `send an empty datagram to "127.0.0.1:9000"`, send.AnEmptyDatagram().To(address).String())
		assert.Equal(
			t,
			`send the datagram "pong" to the sender of the last datagram`,
			send.TheDatagram("pong").ToTheSenderOfTheLastDatagram().String(),
		)
		assert.Equal(
			t,
			`send the datagram "pong" to "127.0.0.1:9000"`,
			send.TheDatagram("pong").ToTheSenderOfTheLastDatagram().To(address).String(),
		)
	})
}

// readFrom reads the next datagram of the connection and tells who sent it.
func readFrom(t *testing.T, conn *net.UDPConn) ([]byte, *net.UDPAddr) {
	t.Helper()

	require.NoError(t, conn.SetReadDeadline(time.Now().Add(timeout)))
	buffer := make([]byte, 1024)
	read, sender, err := conn.ReadFromUDP(buffer)
	require.NoError(t, err)

	return buffer[:read], sender
}
