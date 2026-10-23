package text_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
	"github.com/grandper/go-screenplay-udp/fixture"
	"github.com/grandper/go-screenplay-udp/question"
	questionfixture "github.com/grandper/go-screenplay-udp/question/fixture"
	"github.com/grandper/go-screenplay-udp/question/text"
)

func TestOf(t *testing.T) {
	t.Run("answers with the text of the datagram", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
		require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping").To(echo.String()), receive.TheNextDatagram()))

		err := adam.Should(see.The(text.Of(last.Of(question.UDPDatagrams()))).Should(be.EqualTo("ping")))

		require.NoError(t, err)
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the datagram is nil", func(t *testing.T) {
		wrapped := screenplayfixture.NewFakeQuestion[*ability.Datagram]("the datagram", nil)

		_, err := screenplay.ActorNamed("Adam").AsksFor(text.Of(wrapped))

		require.ErrorIs(t, err, question.ErrDatagramIsNil)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		pong := questionfixture.CreateFakeDatagram("pong", "127.0.0.1:9000")
		wrapped := screenplayfixture.NewFakeQuestion("the datagram", pong)
		assert.Equal(t, "text of the datagram", text.Of(wrapped).String())
		assert.Equal(t, "text of last of UDP datagrams", text.Of(last.Of(question.UDPDatagrams())).String())
	})
}
