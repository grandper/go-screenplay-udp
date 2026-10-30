package local_test

import (
	"net"
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability/communicate"
	"github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/question/local"
	"github.com/grandper/go-screenplay-udp/udp"
)

func TestUDPAddress(t *testing.T) {
	t.Run("answers with the address the connection is bound to", func(t *testing.T) {
		conn := fixture.ListenOnAFreeUDPPort(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().UsingTheConnection(conn))

		address, err := adam.AsksFor(local.UDPAddress())

		require.NoError(t, err)
		assert.Equal(t, conn.LocalAddr().String(), address.String())
	})

	t.Run("answers with a copy", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())

		address, err := adam.AsksFor(local.UDPAddress())
		require.NoError(t, err)
		address.Port = 0

		again, err := adam.AsksFor(local.UDPAddress())
		require.NoError(t, err)
		assert.Positive(t, again.Port)
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the connection is not set", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
		_, err := adam.AsksFor(local.UDPAddress())
		require.ErrorIs(t, err, udp.ErrUDPConnectionNotSet)
	})

	t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
		_, err := screenplay.ActorNamed("Adam").AsksFor(local.UDPAddress())
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "local UDP address", local.UDPAddress().String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[*net.UDPAddr])(nil), local.UDPAddress())
	})
}
