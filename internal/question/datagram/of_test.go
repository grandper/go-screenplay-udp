package datagram_test

import (
	"errors"
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/question/fixture"
	"github.com/grandper/go-screenplay-udp/question/internal/datagram"
	"github.com/grandper/go-screenplay-udp/udp"
)

func TestOf(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")
	pong := fixture.CreateFakeDatagram("pong", "127.0.0.1:9000")

	t.Run("answers with the part of the datagram read", func(t *testing.T) {
		wrapped := screenplayfixture.NewFakeQuestion("the datagram", pong)

		text, err := adam.AsksFor(datagram.Of("text of", wrapped, (*udp.Datagram).Text))

		require.NoError(t, err)
		assert.Equal(t, "pong", text)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		errFailed := errors.New("failed")
		wrapped := screenplayfixture.NewFailingFakeQuestion[*udp.Datagram]("the datagram", errFailed)

		_, err := adam.AsksFor(datagram.Of("text of", wrapped, (*udp.Datagram).Text))

		require.ErrorIs(t, err, errFailed)
	})

	t.Run("fails when the datagram is nil", func(t *testing.T) {
		wrapped := screenplayfixture.NewFakeQuestion[*udp.Datagram]("the datagram", nil)

		_, err := adam.AsksFor(datagram.Of("text of", wrapped, (*udp.Datagram).Text))

		require.ErrorIs(t, err, udp.ErrDatagramIsNil)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		wrapped := screenplayfixture.NewFakeQuestion("the datagram", pong)
		assert.Equal(t, "text of the datagram", datagram.Of("text of", wrapped, (*udp.Datagram).Text).String())
	})
}
