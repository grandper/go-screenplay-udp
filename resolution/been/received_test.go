package been_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/resolution/fixture"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
	udpfixture "github.com/grandper/go-screenplay-udp/fixture"
	internalfixture "github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/question"
	questionfixture "github.com/grandper/go-screenplay-udp/question/fixture"
	"github.com/grandper/go-screenplay-udp/resolution/been"
)

func TestReceived(t *testing.T) {
	ping := questionfixture.CreateFakeDatagram("ping", "127.0.0.1:9000")
	pong := questionfixture.CreateFakeDatagram("pong", "127.0.0.1:9001")
	unknown := questionfixture.CreateFakeDatagram("?", "not an address")

	t.Run("matches when at least one datagram has been received", func(t *testing.T) {
		fixture.AssertMatch(t, been.Received(), []*ability.Datagram{ping})
		fixture.AssertMatch(t, been.Received(), []*ability.Datagram{ping, pong})
	})

	t.Run("does not match when no datagram has been received", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.Received(), []*ability.Datagram{})
		fixture.AssertNoMatch(t, been.Received(), nil)
	})

	t.Run("matches when a datagram has been received from the address", func(t *testing.T) {
		fixture.AssertMatch(t, been.ReceivedFrom("127.0.0.1:9001"), []*ability.Datagram{ping, pong})
		fixture.AssertMatch(t, been.ReceivedFrom("localhost:9001"), []*ability.Datagram{ping, pong})
	})

	t.Run("does not match when no datagram has been received from the address", func(t *testing.T) {
		fixture.AssertNoMatch(t, been.ReceivedFrom("127.0.0.1:9001"), []*ability.Datagram{ping, unknown})
		fixture.AssertNoMatch(t, been.ReceivedFrom("127.0.0.1:9001"), []*ability.Datagram{nil})
		fixture.AssertNoMatch(t, been.ReceivedFrom("127.0.0.1:9001"), nil)
	})

	t.Run("is invalid when the address cannot be resolved", func(t *testing.T) {
		fixture.AssertInvalid(t, been.ReceivedFrom("not an address"), []*ability.Datagram{ping})
	})

	t.Run("reads as a sentence with the question", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		echo := udpfixture.RunAnEchoUDPServer(t)
		adam := internalfixture.NarratedActor("Adam", recorder, udpfixture.ListenOnAFreeUDPPort(t))

		err := adam.Should(see.The(question.UDPDatagramsFrom(echo.String())).ShouldHave(been.Received()))
		require.EqualError(
			t,
			err,
			`Adam failed to see that the UDP datagrams from "`+echo.String()+
				`" should have been received: no datagram has been received`,
		)

		require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping").To(echo.String()), receive.TheNextDatagram()))
		require.NoError(t, adam.Should(see.The(question.UDPDatagrams()).ShouldHave(been.ReceivedFrom(echo.String()))))
		assert.Contains(
			t,
			recorder.Messages(),
			`Adam should see that the UDP datagrams should have been received from "`+echo.String()+`"`,
		)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "been received", been.Received().String())
		assert.Equal(t, `been received from "127.0.0.1:9000"`, been.ReceivedFrom("127.0.0.1:9000").String())
	})
}
