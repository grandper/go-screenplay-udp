package question_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/question"
	"github.com/grandper/go-screenplay-udp/question/payload"
	"github.com/grandper/go-screenplay-udp/question/sender"
	"github.com/grandper/go-screenplay-udp/question/text"
)

func TestErrDatagramIsNil(t *testing.T) {
	t.Run("describes the mistake", func(t *testing.T) {
		require.EqualError(t, question.ErrDatagramIsNil, "the datagram is nil")
	})

	t.Run("is wrapped by the error of a question asked about a nil datagram", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		noDatagram := screenplayfixture.NewFakeQuestion[*ability.Datagram]("the datagram", nil)

		_, err := adam.AsksFor(payload.Of(noDatagram))
		require.ErrorIs(t, err, question.ErrDatagramIsNil)
		_, err = adam.AsksFor(text.Of(noDatagram))
		require.ErrorIs(t, err, question.ErrDatagramIsNil)
		_, err = adam.AsksFor(sender.Of(noDatagram))
		require.ErrorIs(t, err, question.ErrDatagramIsNil)
	})
}
