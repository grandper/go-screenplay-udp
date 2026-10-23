package send_test

import (
	"net"
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/action"
	"github.com/grandper/go-screenplay-udp/action/connect"
	impl "github.com/grandper/go-screenplay-udp/action/internal/send"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/fixture"
	"github.com/grandper/go-screenplay-udp/question"
)

const (
	address = "127.0.0.1:9000"
	text    = "Hello, UDP!"
	timeout = 2 * time.Second
)

func TestTheDatagram(t *testing.T) {
	t.Run("sends the text as the payload", func(t *testing.T) {
		peer := fixture.ListenOnAFreeUDPPort(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())

		require.NoError(t, adam.AttemptsTo(impl.TheDatagram(text).To(peer.LocalAddr().String())))

		data, sender := readFrom(t, peer)
		assert.Equal(t, []byte(text), data)
		local, err := adam.AsksFor(question.LocalUDPAddress())
		require.NoError(t, err)
		assert.Equal(t, local.String(), sender.String())
		require.NoError(t, adam.Exit())
	})

	t.Run("goes to the address the actor connected to without a destination", func(t *testing.T) {
		peer := fixture.ListenOnAFreeUDPPort(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP())
		require.NoError(t, adam.AttemptsTo(connect.ToTheUDPAddress(peer.LocalAddr().String())))

		require.NoError(t, adam.AttemptsTo(impl.TheDatagram(text)))

		data, _ := readFrom(t, peer)
		assert.Equal(t, []byte(text), data)
		require.NoError(t, adam.Exit())
	})

	t.Run("fails without a destination when the actor is not connected", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
		require.ErrorIs(t, adam.AttemptsTo(impl.TheDatagram(text)), ability.ErrDestinationNotSet)
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
		err := screenplay.ActorNamed("Adam").AttemptsTo(impl.TheDatagram(text).To(address))
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `send the datagram "Hello, UDP!"`, impl.TheDatagram(text).String())
		assert.Equal(t, `send the datagram ""`, impl.TheDatagram("").String())
		assert.Equal(
			t,
			`send the datagram "Hello, UDP!" to "127.0.0.1:9000"`,
			impl.TheDatagram(text).To(address).String(),
		)
	})

	t.Run("implements the performable interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Performable)(nil), impl.TheDatagram(text))
	})
}

func TestTheBinaryDatagram(t *testing.T) {
	binary := []byte{0x00, 0x01, 0x02}

	t.Run("sends the bytes as the payload", func(t *testing.T) {
		peer := fixture.ListenOnAFreeUDPPort(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())

		require.NoError(t, adam.AttemptsTo(impl.TheBinaryDatagram(binary).To(peer.LocalAddr().String())))

		data, _ := readFrom(t, peer)
		assert.Equal(t, binary, data)
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the datagram is too large", func(t *testing.T) {
		peer := fixture.ListenOnAFreeUDPPort(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())

		err := adam.AttemptsTo(impl.TheBinaryDatagram(make([]byte, 70000)).To(peer.LocalAddr().String()))

		require.ErrorIs(t, err, ability.ErrFailedToSend)
		require.NoError(t, adam.Exit())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, `send the binary datagram "\x00\x01\x02"`, impl.TheBinaryDatagram(binary).String())
		assert.Equal(
			t,
			`send the binary datagram "\x00\x01\x02" to "127.0.0.1:9000"`,
			impl.TheBinaryDatagram(binary).To(address).String(),
		)
	})
}

func TestAnEmptyDatagram(t *testing.T) {
	t.Run("sends a datagram with no payload", func(t *testing.T) {
		peer := fixture.ListenOnAFreeUDPPort(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())

		require.NoError(t, adam.AttemptsTo(impl.AnEmptyDatagram().To(peer.LocalAddr().String())))

		data, _ := readFrom(t, peer)
		assert.Empty(t, data)
		require.NoError(t, adam.Exit())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "send an empty datagram", impl.AnEmptyDatagram().String())
		assert.Equal(t, `send an empty datagram to "127.0.0.1:9000"`, impl.AnEmptyDatagram().To(address).String())
	})
}

func TestTheDatagramAction_To(t *testing.T) {
	t.Run("fails when the address is empty", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
		require.ErrorIs(t, adam.AttemptsTo(impl.TheDatagram(text).To("")), action.ErrAddressIsEmpty)
		require.NoError(t, adam.Exit())
	})

	t.Run("reports an empty address before looking for the ability", func(t *testing.T) {
		err := screenplay.ActorNamed("Adam").AttemptsTo(impl.TheDatagram(text).To(""))
		require.ErrorIs(t, err, action.ErrAddressIsEmpty)
		require.NotErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("fails when the address cannot be resolved", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
		require.ErrorIs(t, adam.AttemptsTo(impl.TheDatagram(text).To("not an address")), ability.ErrInvalidAddress)
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the actor is connected to another address", func(t *testing.T) {
		peer := fixture.ListenOnAFreeUDPPort(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP())
		require.NoError(t, adam.AttemptsTo(connect.ToTheUDPAddress(peer.LocalAddr().String())))

		err := adam.AttemptsTo(impl.TheDatagram(text).To(address))

		require.ErrorIs(t, err, ability.ErrConnectedToAnotherAddress)
		require.NoError(t, adam.Exit())
	})

	t.Run("keeps the last destination written", func(t *testing.T) {
		theAction := impl.TheDatagram(text).To("127.0.0.1:9001").To(address)
		assert.Equal(t, `send the datagram "Hello, UDP!" to "127.0.0.1:9000"`, theAction.String())

		theAction = impl.TheDatagram(text).ToTheSenderOfTheLastDatagram().To(address)
		assert.Equal(t, `send the datagram "Hello, UDP!" to "127.0.0.1:9000"`, theAction.String())
	})
}

func TestTheDatagramAction_ToTheSenderOfTheLastDatagram(t *testing.T) {
	t.Run("replies to the sender of the last received datagram", func(t *testing.T) {
		first := fixture.ListenOnAFreeUDPPort(t)
		last := fixture.ListenOnAFreeUDPPort(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
		local, err := adam.AsksFor(question.LocalUDPAddress())
		require.NoError(t, err)
		_, err = first.WriteToUDP([]byte("ping"), local)
		require.NoError(t, err)
		require.NoError(t, adam.AttemptsTo(receive.TheNextDatagram()))
		_, err = last.WriteToUDP([]byte("ping"), local)
		require.NoError(t, err)
		require.NoError(t, adam.AttemptsTo(receive.TheNextDatagram()))

		require.NoError(t, adam.AttemptsTo(impl.TheDatagram("pong").ToTheSenderOfTheLastDatagram()))

		data, sender := readFrom(t, last)
		assert.Equal(t, []byte("pong"), data)
		assert.Equal(t, local.String(), sender.String())
		require.NoError(t, first.SetReadDeadline(time.Now().Add(100*time.Millisecond)))
		_, _, err = first.ReadFromUDP(make([]byte, 1024))
		require.Error(t, err)
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when no datagram has been received", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
		err := adam.AttemptsTo(impl.TheDatagram("pong").ToTheSenderOfTheLastDatagram())
		require.ErrorIs(t, err, ability.ErrNoDatagramReceived)
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
		err := screenplay.ActorNamed("Adam").AttemptsTo(impl.TheDatagram("pong").ToTheSenderOfTheLastDatagram())
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("replaces the address written before", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
		theAction := impl.TheDatagram("pong").To("").ToTheSenderOfTheLastDatagram()

		assert.Equal(t, `send the datagram "pong" to the sender of the last datagram`, theAction.String())
		require.ErrorIs(t, adam.AttemptsTo(theAction), ability.ErrNoDatagramReceived)
		require.NoError(t, adam.Exit())
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
