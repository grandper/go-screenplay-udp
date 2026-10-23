package question_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/fixture"
	"github.com/grandper/go-screenplay-udp/question"
)

func TestLocalUDPAddress(t *testing.T) {
	t.Run("answers with the address the connection is bound to", func(t *testing.T) {
		conn := fixture.ListenOnAFreeUDPPort(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().UsingTheConnection(conn))

		local, err := adam.AsksFor(question.LocalUDPAddress())

		require.NoError(t, err)
		assert.Equal(t, conn.LocalAddr().String(), local.String())
	})

	t.Run("answers with a copy", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())

		local, err := adam.AsksFor(question.LocalUDPAddress())
		require.NoError(t, err)
		local.Port = 0

		again, err := adam.AsksFor(question.LocalUDPAddress())
		require.NoError(t, err)
		assert.Positive(t, again.Port)
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the connection is not set", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP())
		_, err := adam.AsksFor(question.LocalUDPAddress())
		require.ErrorIs(t, err, ability.ErrUDPConnectionNotSet)
	})

	t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
		_, err := screenplay.ActorNamed("Adam").AsksFor(question.LocalUDPAddress())
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "local UDP address", question.LocalUDPAddress().String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[*net.UDPAddr])(nil), question.LocalUDPAddress())
	})
}

// portOf returns the port of the address as a text.
func portOf(address *net.UDPAddr) string {
	return strconv.Itoa(address.Port)
}
