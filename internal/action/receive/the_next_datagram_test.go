package receive_test

import (
	"context"
	"testing"
	"time"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/question/number"
	"github.com/grandper/go-screenplay/v2/resolution/be"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability/communicate"
	"github.com/grandper/go-screenplay-udp/action/send"
	impl "github.com/grandper/go-screenplay-udp/internal/action/receive"
	"github.com/grandper/go-screenplay-udp/internal/fixture"
	"github.com/grandper/go-screenplay-udp/udp"
)

func TestTheNextDatagram(t *testing.T) {
	t.Run("receives one datagram and stores it", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
		require.NoError(
			t,
			adam.AttemptsTo(send.TheDatagram("1").To(echo.String()), send.TheDatagram("2").To(echo.String())),
		)

		require.NoError(t, adam.AttemptsTo(impl.TheNextDatagram()))

		datagrams, err := adam.AsksFor(udp.Datagrams())
		require.NoError(t, err)
		require.Len(t, datagrams, 1)
		assert.Equal(t, "1", datagrams[0].Text())
		assert.Equal(t, echo.String(), datagrams[0].Sender().String())
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the connection is not set", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())
		require.ErrorIs(t, adam.AttemptsTo(impl.TheNextDatagram()), udp.ErrUDPConnectionNotSet)
	})

	t.Run("fails when the actor cannot communicate over UDP", func(t *testing.T) {
		err := screenplay.ActorNamed("Adam").AttemptsTo(impl.TheNextDatagram())
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "receive the next datagram", impl.TheNextDatagram().String())
	})

	t.Run("implements the performable interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Performable)(nil), impl.TheNextDatagram())
	})
}

func TestTheNext(t *testing.T) {
	t.Run("receives as many datagrams as the amount, in order", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
		require.NoError(t, adam.AttemptsTo(
			send.TheDatagram("1").To(echo.String()),
			send.TheDatagram("2").To(echo.String()),
			send.TheDatagram("3").To(echo.String()),
		))

		require.NoError(t, adam.AttemptsTo(impl.TheNext(2).Datagrams()))

		datagrams, err := adam.AsksFor(udp.Datagrams())
		require.NoError(t, err)
		require.Len(t, datagrams, 2)
		assert.Equal(t, "1", datagrams[0].Text())
		assert.Equal(t, "2", datagrams[1].Text())
		require.NoError(t, adam.Exit())
	})

	t.Run("receives one datagram in the singular", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
		require.NoError(t, adam.AttemptsTo(send.TheDatagram("1").To(echo.String())))

		require.NoError(t, adam.AttemptsTo(impl.TheNext(1).Datagram()))

		require.NoError(t, adam.Should(see.The(number.Of(udp.Datagrams())).Should(be.EqualTo(1))))
		require.NoError(t, adam.Exit())
	})

	t.Run("keeps the datagrams received before the one that fails", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
		require.NoError(t, adam.AttemptsTo(send.TheDatagram("1").To(echo.String())))

		err := adam.AttemptsTo(impl.TheNext(2).Datagrams().TimingOutAfter(50).Milliseconds())

		require.ErrorIs(t, err, udp.ErrFailedToReceive)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.NoError(t, adam.Should(see.The(number.Of(udp.Datagrams())).Should(be.EqualTo(1))))
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the datagrams kept would exceed the buffered size of the actor", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.NewProduction(screenplay.WithMaxBufferedSize(1)).
			SetTheStage(screenplay.CastOfStandardActors()).
			ActorNamed("Adam").
			WhoCan(communicate.OverUDP().ListeningOnAFreePort())
		require.NoError(
			t,
			adam.AttemptsTo(send.TheDatagram("1").To(echo.String()), send.TheDatagram("2").To(echo.String())),
		)

		err := adam.AttemptsTo(impl.TheNext(2).Datagrams())

		require.ErrorIs(t, err, udp.ErrDatagramsExceedTheBound)
		require.NoError(t, adam.Should(see.The(number.Of(udp.Datagrams())).Should(be.EqualTo(1))))
		require.NoError(t, adam.Exit())
	})

	t.Run("fails when the amount is not positive", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
		require.ErrorIs(t, adam.AttemptsTo(impl.TheNext(0).Datagrams()), udp.ErrNumberOfDatagramsIsNotPositive)
		require.ErrorIs(t, adam.AttemptsTo(impl.TheNext(-1).Datagram()), udp.ErrNumberOfDatagramsIsNotPositive)
		require.NoError(t, adam.Exit())
	})

	t.Run("reports an amount that is not positive before looking for the ability", func(t *testing.T) {
		err := screenplay.ActorNamed("Adam").AttemptsTo(impl.TheNext(0).Datagrams())
		require.ErrorIs(t, err, udp.ErrNumberOfDatagramsIsNotPositive)
		require.NotErrorIs(t, err, screenplay.ErrActorMissingAbility)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "receive the next 3 datagrams", impl.TheNext(3).Datagrams().String())
		assert.Equal(t, "receive the next 1 datagram", impl.TheNext(1).Datagram().String())
		assert.Equal(t, "receive the next 1 datagram", impl.TheNext(1).Datagrams().String())
		assert.Equal(t, "receive the next 2 datagrams", impl.TheNext(2).Datagram().String())
		assert.Equal(t, "receive the next 0 datagrams", impl.TheNext(0).Datagrams().String())
	})
}

func TestTheNextDatagramAction_TimingOutAfter(t *testing.T) {
	t.Run("waits at most the timeout of the actor when it is not written", func(t *testing.T) {
		production := screenplay.NewProduction(screenplay.WithTimeout(50 * time.Millisecond))
		adam := production.SetTheStage(screenplay.CastOfStandardActors()).
			ActorNamed("Adam").
			WhoCan(communicate.OverUDP().ListeningOnAFreePort())

		start := time.Now()
		err := adam.AttemptsTo(impl.TheNextDatagram())

		require.ErrorIs(t, err, udp.ErrFailedToReceive)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Less(t, time.Since(start), time.Second)
		require.NoError(t, adam.Exit())
	})

	t.Run("waits at most the timeout written instead of the one of the actor", func(t *testing.T) {
		production := screenplay.NewProduction(screenplay.WithTimeout(time.Hour))
		adam := production.SetTheStage(screenplay.CastOfStandardActors()).
			ActorNamed("Adam").
			WhoCan(communicate.OverUDP().ListeningOnAFreePort())

		start := time.Now()
		err := adam.AttemptsTo(impl.TheNextDatagram().TimingOutAfter(50).Milliseconds())

		require.ErrorIs(t, err, udp.ErrFailedToReceive)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Less(t, time.Since(start), time.Second)
		require.NoError(t, adam.Exit())
	})

	t.Run("leaves time for a datagram that comes before the timeout", func(t *testing.T) {
		echo := fixture.RunAnEchoUDPServer(t)
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
		require.NoError(t, adam.AttemptsTo(send.TheDatagram("ping").To(echo.String())))

		require.NoError(t, adam.AttemptsTo(impl.TheNextDatagram().TimingOutAfter(2).Seconds()))

		require.NoError(t, adam.Should(see.The(number.Of(udp.Datagrams())).Should(be.EqualTo(1))))
		require.NoError(t, adam.Exit())
	})

	t.Run("stops when the context of the actor ends", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		adam := screenplay.ActorNamed("Adam").
			WithContext(ctx).
			WhoCan(communicate.OverUDP().ListeningOnAFreePort())
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		err := adam.AttemptsTo(impl.TheNextDatagram().TimingOutAfter(1).Minute())

		require.ErrorIs(t, err, udp.ErrFailedToReceive)
		require.ErrorIs(t, err, context.Canceled)
		require.NoError(t, adam.Exit())
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(
			t,
			"receive the next datagram timing out after 2 seconds",
			impl.TheNextDatagram().TimingOutAfter(2).Seconds().String(),
		)
		assert.Equal(
			t,
			"receive the next 3 datagrams timing out after 500 milliseconds",
			impl.TheNext(3).Datagrams().TimingOutAfter(500).Milliseconds().String(),
		)
		assert.Equal(
			t,
			"receive the next 1 datagram timing out after 1 minute",
			impl.TheNext(1).Datagram().TimingOutAfter(1).Minute().String(),
		)
	})
}
