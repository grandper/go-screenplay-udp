package receive_test

import (
	"context"
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/question/number"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/action"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
	"github.com/grandper/go-screenplay-udp/fixture"
	internalfixture "github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/question"
)

func TestTheNextDatagram(t *testing.T) {
	t.Run("receives the next datagram", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			echo := fixture.RunAnEchoUDPServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping").To(echo.String())))

			require.NoError(t, adam.AttemptsTo(receive.TheNextDatagram()))

			datagrams, err := adam.AsksFor(question.UDPDatagrams())
			require.NoError(t, err)
			require.Len(t, datagrams, 1)
			assert.Equal(t, "ping", datagrams[0].Text())
			assert.Equal(t, echo.String(), datagrams[0].Sender().String())
			require.NoError(t, adam.Exit())
		})

		t.Run("several at once", func(t *testing.T) {
			echo := fixture.RunAnEchoUDPServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			require.NoError(t, adam.AttemptsTo(
				send.TheDatagram("1").To(echo.String()),
				send.TheDatagram("2").To(echo.String()),
				send.TheDatagram("3").To(echo.String()),
			))

			require.NoError(t, adam.AttemptsTo(receive.TheNext(2).Datagrams(), receive.TheNext(1).Datagram()))

			require.NoError(t, adam.Should(see.The(number.Of(question.UDPDatagrams())).Should(be.EqualTo(3))))
			require.NoError(t, adam.Exit())
		})

		t.Run("waits at most the timeout of the actor", func(t *testing.T) {
			production := screenplay.NewProduction(screenplay.WithTimeout(50 * time.Millisecond))
			adam := production.SetTheStage(screenplay.CastOfStandardActors()).
				ActorNamed("Adam").
				WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())

			start := time.Now()
			err := adam.AttemptsTo(receive.TheNextDatagram())

			require.ErrorIs(t, err, ability.ErrFailedToReceive)
			require.ErrorIs(t, err, context.DeadlineExceeded)
			assert.Less(t, time.Since(start), time.Second)
			require.NoError(t, adam.Exit())
		})

		t.Run("waits at most the timeout written", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())

			start := time.Now()
			err := adam.AttemptsTo(receive.TheNextDatagram().TimingOutAfter(50).Milliseconds())

			require.ErrorIs(t, err, ability.ErrFailedToReceive)
			require.ErrorIs(t, err, context.DeadlineExceeded)
			assert.Less(t, time.Since(start), time.Second)
			require.NoError(t, adam.Exit())
		})

		t.Run("waits at most the timeout for all of them", func(t *testing.T) {
			echo := fixture.RunAnEchoUDPServer(t)
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			require.NoError(t, adam.AttemptsTo(send.TheDatagram("1").To(echo.String())))

			start := time.Now()
			err := adam.AttemptsTo(receive.TheNext(2).Datagrams().TimingOutAfter(50).Milliseconds())

			require.ErrorIs(t, err, context.DeadlineExceeded)
			assert.Less(t, time.Since(start), time.Second)
			require.NoError(t, adam.Should(see.The(number.Of(question.UDPDatagrams())).Should(be.EqualTo(1))))
			require.NoError(t, adam.Exit())
		})

		t.Run("stops when the context of the actor ends", func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			adam := screenplay.ActorNamed("Adam").
				WithContext(ctx).
				WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			go func() {
				time.Sleep(50 * time.Millisecond)
				cancel()
			}()

			err := adam.AttemptsTo(receive.TheNextDatagram())

			require.ErrorIs(t, err, context.Canceled)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the number of datagrams is not positive", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
			require.ErrorIs(
				t,
				adam.AttemptsTo(receive.TheNext(0).Datagrams()),
				action.ErrNumberOfDatagramsIsNotPositive,
			)
			require.ErrorIs(
				t,
				adam.AttemptsTo(receive.TheNext(-1).Datagrams()),
				action.ErrNumberOfDatagramsIsNotPositive,
			)
			require.NoError(t, adam.Exit())
		})

		t.Run("fails when the connection is not set", func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP())
			require.ErrorIs(t, adam.AttemptsTo(receive.TheNextDatagram()), ability.ErrUDPConnectionNotSet)
		})

		t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
			err := screenplay.ActorNamed("Adam").AttemptsTo(receive.TheNextDatagram())
			require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		})
	})

	t.Run("is narrated", func(t *testing.T) {
		recorder := screenplayfixture.NewRecorder()
		echo := fixture.RunAnEchoUDPServer(t)
		adam := internalfixture.NarratedActor("Adam", recorder, fixture.ListenOnAFreeUDPPort(t))
		require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping").To(echo.String())))

		require.NoError(t, adam.AttemptsTo(receive.TheNextDatagram()))

		assert.Contains(t, recorder.Messages(), "Adam attempts to receive the next datagram")
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "receive the next datagram", receive.TheNextDatagram().String())
		assert.Equal(t, "receive the next 3 datagrams", receive.TheNext(3).Datagrams().String())
		assert.Equal(t, "receive the next 1 datagram", receive.TheNext(1).Datagram().String())
		assert.Equal(t, "receive the next 1 datagram", receive.TheNext(1).Datagrams().String())
		assert.Equal(
			t,
			"receive the next datagram timing out after 2 seconds",
			receive.TheNextDatagram().TimingOutAfter(2).Seconds().String(),
		)
		assert.Equal(
			t,
			"receive the next 3 datagrams timing out after 500 milliseconds",
			receive.TheNext(3).Datagrams().TimingOutAfter(500).Milliseconds().String(),
		)
	})
}
