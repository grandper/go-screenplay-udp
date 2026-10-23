package ability_test

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/fixture"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
)

const (
	timeout = 2 * time.Second
	// unroutableAddress is an address of TEST-NET-1 (RFC 5737), which no
	// interface of the host has, so it cannot be bound.
	unroutableAddress = "192.0.2.1:0"
)

func TestCommunicateOverUDP(t *testing.T) {
	t.Run("is given a connection", func(t *testing.T) {
		t.Run("and keeps it", func(t *testing.T) {
			conn := fixture.ListenOnAFreeUDPPort(t)
			other := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").
				WhoCan(impl.CommunicateOverUDP().UsingTheConnection(conn).UsingTheConnection(other))

			local, err := impl.TheAbilityOf(adam).ToTellTheLocalAddress()

			require.NoError(t, err)
			assert.Equal(t, conn.LocalAddr().String(), local.String())
		})

		t.Run("and ignores a nil one", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().UsingTheConnection(nil))

			_, err := impl.TheAbilityOf(adam).ToTellTheLocalAddress()

			require.ErrorIs(t, err, impl.ErrUDPConnectionNotSet)
		})

		t.Run("on a nil receiver", func(t *testing.T) {
			var ability *impl.CommunicateOverUDPAbility
			assert.Nil(t, ability.UsingTheConnection(fixture.ListenOnAFreeUDPPort(t)))
			assert.Nil(t, ability.ListeningOn("127.0.0.1:0"))
			assert.Nil(t, ability.ListeningOnAFreePort())
		})
	})

	t.Run("is given an address to listen on", func(t *testing.T) {
		t.Run("and binds it the first time it is used", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOn("127.0.0.1:0"))

			local, err := impl.TheAbilityOf(adam).ToTellTheLocalAddress()

			require.NoError(t, err)
			assert.Equal(t, "127.0.0.1", local.IP.String())
			assert.Positive(t, local.Port)
			require.NoError(t, adam.Exit())
		})

		t.Run("or a free port", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())

			local, err := impl.TheAbilityOf(adam).ToTellTheLocalAddress()

			require.NoError(t, err)
			assert.Equal(t, "127.0.0.1", local.IP.String())
			assert.Positive(t, local.Port)
			require.NoError(t, adam.Exit())
		})

		t.Run("and reports an address it cannot resolve", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOn("not an address"))

			_, err := impl.TheAbilityOf(adam).ToTellTheLocalAddress()

			require.ErrorIs(t, err, impl.ErrInvalidAddress)
		})

		t.Run("and reports an address it cannot bind", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOn(unroutableAddress))

			_, err := impl.TheAbilityOf(adam).ToTellTheLocalAddress()

			require.ErrorIs(t, err, impl.ErrFailedToListen)
		})
	})

	t.Run("is given nothing", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP())

		require.ErrorIs(
			t,
			impl.TheAbilityOf(adam).ToSend([]byte("ping"), "127.0.0.1:9000"),
			impl.ErrUDPConnectionNotSet,
		)
		require.ErrorIs(t, impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(t.Context()), impl.ErrUDPConnectionNotSet)
		require.ErrorIs(t, impl.TheAbilityOf(adam).ToClose(), impl.ErrUDPConnectionNotSet)
		_, err := impl.TheAbilityOf(adam).ToRetrieveDatagrams()
		require.ErrorIs(t, err, impl.ErrUDPConnectionNotSet)
	})

	t.Run("listens", func(t *testing.T) {
		t.Run("on an address", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP())

			require.NoError(t, impl.TheAbilityOf(adam).ToListenOn("127.0.0.1:0"))

			local, err := impl.TheAbilityOf(adam).ToTellTheLocalAddress()
			require.NoError(t, err)
			assert.Positive(t, local.Port)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when it already holds a connection", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())
			require.NoError(t, impl.TheAbilityOf(adam).ToListenOn("127.0.0.1:0"))

			require.ErrorIs(t, impl.TheAbilityOf(adam).ToListenOn("127.0.0.1:0"), impl.ErrAlreadyListening)
			require.ErrorIs(t, impl.TheAbilityOf(adam).ToConnectTo("127.0.0.1:9000"), impl.ErrAlreadyListening)
			require.ErrorIs(
				t,
				impl.TheAbilityOf(adam).ToListenOnTheMulticastGroup("224.0.0.251:0"),
				impl.ErrAlreadyListening,
			)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the address cannot be resolved or bound", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP())

			require.ErrorIs(t, impl.TheAbilityOf(adam).ToListenOn("not an address"), impl.ErrInvalidAddress)
			require.ErrorIs(t, impl.TheAbilityOf(adam).ToListenOn(unroutableAddress), impl.ErrFailedToListen)
		})

		t.Run("on a multicast group", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP())

			require.NoError(t, impl.TheAbilityOf(adam).ToListenOnTheMulticastGroup("224.0.0.251:0"))

			local, err := impl.TheAbilityOf(adam).ToTellTheLocalAddress()
			require.NoError(t, err)
			assert.Positive(t, local.Port)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the multicast group cannot be resolved or joined", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP())

			require.ErrorIs(
				t,
				impl.TheAbilityOf(adam).ToListenOnTheMulticastGroup("not an address"),
				impl.ErrInvalidAddress,
			)
			require.ErrorIs(
				t,
				impl.TheAbilityOf(adam).ToListenOnTheMulticastGroup("127.0.0.1:0"),
				impl.ErrFailedToListen,
			)
		})
	})

	t.Run("connects", func(t *testing.T) {
		t.Run("to a remote address, and sends to it without an address", func(t *testing.T) {
			echo := fixture.RunAnEchoUDPServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP())

			require.NoError(t, impl.TheAbilityOf(adam).ToConnectTo(echo.String()))
			require.NoError(t, impl.TheAbilityOf(adam).ToSend([]byte("ping"), ""))
			require.NoError(t, impl.TheAbilityOf(adam).ToSend([]byte("ping"), echo.String()))

			ctx, cancel := context.WithTimeout(t.Context(), timeout)
			defer cancel()
			require.NoError(t, impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx))
			require.NoError(t, impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx))
			datagrams, err := impl.TheAbilityOf(adam).ToRetrieveDatagramsFrom(echo.String())
			require.NoError(t, err)
			assert.Len(t, datagrams, 2)
			require.NoError(t, adam.Exit())
		})

		t.Run("and refuses to send to another address", func(t *testing.T) {
			echo := fixture.RunAnEchoUDPServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP())
			require.NoError(t, impl.TheAbilityOf(adam).ToConnectTo(echo.String()))

			err := impl.TheAbilityOf(adam).ToSend([]byte("ping"), "127.0.0.1:9")

			require.ErrorIs(t, err, impl.ErrConnectedToAnotherAddress)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the address cannot be resolved or reached", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP())

			require.ErrorIs(t, impl.TheAbilityOf(adam).ToConnectTo("not an address"), impl.ErrInvalidAddress)
			require.ErrorIs(t, impl.TheAbilityOf(adam).ToConnectTo("127.0.0.1:0"), impl.ErrFailedToConnect)
		})
	})

	t.Run("sends", func(t *testing.T) {
		t.Run("a datagram to an address", func(t *testing.T) {
			peer := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())

			require.NoError(t, impl.TheAbilityOf(adam).ToSend([]byte("ping"), peer.LocalAddr().String()))

			assert.Equal(t, "ping", readFrom(t, peer))
			require.NoError(t, adam.Exit())
		})

		t.Run("fails without a destination", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())

			require.ErrorIs(t, impl.TheAbilityOf(adam).ToSend([]byte("ping"), ""), impl.ErrDestinationNotSet)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the address cannot be resolved", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())

			require.ErrorIs(t, impl.TheAbilityOf(adam).ToSend([]byte("ping"), "not an address"), impl.ErrInvalidAddress)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the datagram is too large", func(t *testing.T) {
			peer := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())

			err := impl.TheAbilityOf(adam).ToSend(make([]byte, 70000), peer.LocalAddr().String())

			require.ErrorIs(t, err, impl.ErrFailedToSend)
			require.NoError(t, adam.Exit())
		})

		t.Run("to the sender of the last datagram", func(t *testing.T) {
			peer := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())
			local, err := impl.TheAbilityOf(adam).ToTellTheLocalAddress()
			require.NoError(t, err)
			_, err = peer.WriteToUDP([]byte("ping"), local)
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(t.Context(), timeout)
			defer cancel()
			require.NoError(t, impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx))

			require.NoError(t, impl.TheAbilityOf(adam).ToSendToTheSenderOfTheLastDatagram([]byte("pong")))

			assert.Equal(t, "pong", readFrom(t, peer))
			require.NoError(t, adam.Exit())
		})

		t.Run("fails to reply when no datagram has been received", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())

			err := impl.TheAbilityOf(adam).ToSendToTheSenderOfTheLastDatagram([]byte("pong"))

			require.ErrorIs(t, err, impl.ErrNoDatagramReceived)
			require.NoError(t, adam.Exit())
		})
	})

	t.Run("receives", func(t *testing.T) {
		t.Run("the next datagram and stores it", func(t *testing.T) {
			echo := fixture.RunAnEchoUDPServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())
			require.NoError(t, impl.TheAbilityOf(adam).ToSend([]byte("ping"), echo.String()))
			ctx, cancel := context.WithTimeout(t.Context(), timeout)
			defer cancel()

			require.NoError(t, impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx))

			datagrams, err := impl.TheAbilityOf(adam).ToRetrieveDatagrams()
			require.NoError(t, err)
			require.Len(t, datagrams, 1)
			assert.Equal(t, "ping", datagrams[0].Text())
			assert.Equal(t, echo.String(), datagrams[0].Sender().String())
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when no datagram comes in time", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())
			ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
			defer cancel()

			start := time.Now()
			err := impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx)

			require.ErrorIs(t, err, impl.ErrFailedToReceive)
			require.ErrorIs(t, err, context.DeadlineExceeded)
			assert.Less(t, time.Since(start), time.Second)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the context is cancelled", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())
			ctx, cancel := context.WithCancel(t.Context())
			go func() {
				time.Sleep(50 * time.Millisecond)
				cancel()
			}()

			err := impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx)

			require.ErrorIs(t, err, impl.ErrFailedToReceive)
			require.ErrorIs(t, err, context.Canceled)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the context already ended", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())
			ctx, cancel := context.WithCancel(t.Context())
			cancel()

			err := impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx)

			require.ErrorIs(t, err, impl.ErrFailedToReceive)
			require.ErrorIs(t, err, context.Canceled)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the connection is closed", func(t *testing.T) {
			conn := fixture.ListenOnAFreeUDPPort(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().UsingTheConnection(conn))
			require.NoError(t, conn.Close())

			err := impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(t.Context())

			require.ErrorIs(t, err, impl.ErrFailedToReceive)
			require.ErrorIs(t, err, net.ErrClosed)
		})
	})

	t.Run("retrieves the datagrams", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		other := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())
		ctx, cancel := context.WithTimeout(t.Context(), timeout)
		defer cancel()
		require.NoError(t, impl.TheAbilityOf(adam).ToSend([]byte("1"), echo.String()))
		require.NoError(t, impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx))
		require.NoError(t, impl.TheAbilityOf(adam).ToSend([]byte("2"), other.String()))
		require.NoError(t, impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx))

		t.Run("all of them, in the order they were received", func(t *testing.T) {
			datagrams, err := impl.TheAbilityOf(adam).ToRetrieveDatagrams()

			require.NoError(t, err)
			require.Len(t, datagrams, 2)
			assert.Equal(t, "1", datagrams[0].Text())
			assert.Equal(t, "2", datagrams[1].Text())
		})

		t.Run("as a copy", func(t *testing.T) {
			datagrams, err := impl.TheAbilityOf(adam).ToRetrieveDatagrams()
			require.NoError(t, err)
			datagrams[0] = nil

			again, err := impl.TheAbilityOf(adam).ToRetrieveDatagrams()
			require.NoError(t, err)
			assert.NotNil(t, again[0])
		})

		t.Run("those from an address", func(t *testing.T) {
			datagrams, err := impl.TheAbilityOf(adam).ToRetrieveDatagramsFrom(other.String())

			require.NoError(t, err)
			require.Len(t, datagrams, 1)
			assert.Equal(t, "2", datagrams[0].Text())
		})

		t.Run("none from an address nothing was received from", func(t *testing.T) {
			datagrams, err := impl.TheAbilityOf(adam).ToRetrieveDatagramsFrom("127.0.0.1:9")

			require.NoError(t, err)
			assert.Empty(t, datagrams)
		})

		t.Run("fails when the address cannot be resolved", func(t *testing.T) {
			_, err := impl.TheAbilityOf(adam).ToRetrieveDatagramsFrom("not an address")

			require.ErrorIs(t, err, impl.ErrInvalidAddress)
		})

		require.NoError(t, adam.Exit())
	})

	t.Run("closes the connection", func(t *testing.T) {
		t.Run("and can listen again, keeping the datagrams", func(t *testing.T) {
			echo := fixture.RunAnEchoUDPServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())
			ctx, cancel := context.WithTimeout(t.Context(), timeout)
			defer cancel()
			require.NoError(t, impl.TheAbilityOf(adam).ToSend([]byte("ping"), echo.String()))
			require.NoError(t, impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx))

			require.NoError(t, impl.TheAbilityOf(adam).ToClose())

			require.ErrorIs(
				t,
				impl.TheAbilityOf(adam).ToSend([]byte("ping"), echo.String()),
				impl.ErrUDPConnectionNotSet,
			)
			require.NoError(t, impl.TheAbilityOf(adam).ToListenOn("127.0.0.1:0"))
			datagrams, err := impl.TheAbilityOf(adam).ToRetrieveDatagrams()
			require.NoError(t, err)
			assert.Len(t, datagrams, 1)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when it holds none", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP())

			require.ErrorIs(t, impl.TheAbilityOf(adam).ToClose(), impl.ErrUDPConnectionNotSet)
		})
	})

	t.Run("forgets", func(t *testing.T) {
		t.Run("the connection and the datagrams", func(t *testing.T) {
			conn := fixture.ListenOnAFreeUDPPort(t)
			echo := fixture.RunAnEchoUDPServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().UsingTheConnection(conn))
			ctx, cancel := context.WithTimeout(t.Context(), timeout)
			defer cancel()
			require.NoError(t, impl.TheAbilityOf(adam).ToSend([]byte("ping"), echo.String()))
			require.NoError(t, impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx))

			require.NoError(t, adam.Exit())

			_, err := conn.WriteToUDP([]byte("ping"), echo)
			require.ErrorIs(t, err, net.ErrClosed)
			adam = screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().UsingTheConnection(conn))
			_, err = impl.TheAbilityOf(adam).ToRetrieveDatagrams()
			require.NoError(t, err)
		})

		t.Run("nothing when it holds nothing", func(t *testing.T) {
			require.NoError(t, impl.CommunicateOverUDP().ListeningOnAFreePort().Forget())
		})

		t.Run("fails on a nil receiver", func(t *testing.T) {
			var ability *impl.CommunicateOverUDPAbility
			require.ErrorIs(t, ability.Forget(), screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("fails every operation when the actor cannot communicate over UDP", func(t *testing.T) {
		for _, actor := range []*screenplay.Actor{nil, screenplay.ActorNamed("Adam")} {
			operations := impl.TheAbilityOf(actor)
			require.ErrorIs(t, operations.ToListenOn("127.0.0.1:0"), screenplay.ErrActorMissingAbility)
			require.ErrorIs(
				t,
				operations.ToListenOnTheMulticastGroup("224.0.0.251:0"),
				screenplay.ErrActorMissingAbility,
			)
			require.ErrorIs(t, operations.ToConnectTo("127.0.0.1:9000"), screenplay.ErrActorMissingAbility)
			require.ErrorIs(t, operations.ToSend([]byte("ping"), "127.0.0.1:9000"), screenplay.ErrActorMissingAbility)
			require.ErrorIs(
				t,
				operations.ToSendToTheSenderOfTheLastDatagram([]byte("pong")),
				screenplay.ErrActorMissingAbility,
			)
			require.ErrorIs(t, operations.ToReceiveTheNextDatagram(t.Context()), screenplay.ErrActorMissingAbility)
			require.ErrorIs(t, operations.ToClose(), screenplay.ErrActorMissingAbility)
			_, err := operations.ToRetrieveDatagrams()
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
			_, err = operations.ToRetrieveDatagramsFrom("127.0.0.1:9000")
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
			_, err = operations.ToTellTheLocalAddress()
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		}
	})

	t.Run("is safe for concurrent use", func(t *testing.T) {
		const senders = 10
		echo := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(impl.CommunicateOverUDP().ListeningOnAFreePort())
		ctx, cancel := context.WithTimeout(t.Context(), timeout)
		defer cancel()

		var group sync.WaitGroup
		for range senders {
			group.Go(func() {
				assert.NoError(t, impl.TheAbilityOf(adam).ToSend([]byte("ping"), echo.String()))
			})
			group.Go(func() {
				assert.NoError(t, impl.TheAbilityOf(adam).ToReceiveTheNextDatagram(ctx))
			})
			group.Go(func() {
				_, err := impl.TheAbilityOf(adam).ToRetrieveDatagrams()
				assert.NoError(t, err)
			})
		}
		group.Wait()

		datagrams, err := impl.TheAbilityOf(adam).ToRetrieveDatagrams()
		require.NoError(t, err)
		assert.Len(t, datagrams, senders)
		require.NoError(t, adam.Exit())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "communicate over UDP", impl.CommunicateOverUDP().String())
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
