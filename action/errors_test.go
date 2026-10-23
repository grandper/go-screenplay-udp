package action_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/action"
	"github.com/grandper/go-screenplay-udp/action/listen"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
)

func TestErrors(t *testing.T) {
	t.Run("describe the mistake", func(t *testing.T) {
		require.EqualError(t, action.ErrAddressIsEmpty, "the UDP address must be provided")
		require.EqualError(
			t,
			action.ErrNumberOfDatagramsIsNotPositive,
			"the number of datagrams to receive must be positive",
		)
	})

	t.Run("are distinct from one another", func(t *testing.T) {
		require.NotErrorIs(t, action.ErrAddressIsEmpty, action.ErrNumberOfDatagramsIsNotPositive)
		require.NotErrorIs(t, action.ErrNumberOfDatagramsIsNotPositive, action.ErrAddressIsEmpty)
	})

	t.Run("are reported when the actor performs the action, not when it is built", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP())

		var sendToNowhere, listenNowhere, receiveNothing screenplay.Performable
		require.NotPanics(t, func() {
			sendToNowhere = send.TheDatagram("ping").To("")
			listenNowhere = listen.OnTheUDPAddress("")
			receiveNothing = receive.TheNext(0).Datagrams()
		})

		require.ErrorIs(t, adam.AttemptsTo(sendToNowhere), action.ErrAddressIsEmpty)
		require.ErrorIs(t, adam.AttemptsTo(listenNowhere), action.ErrAddressIsEmpty)
		require.ErrorIs(t, adam.AttemptsTo(receiveNothing), action.ErrNumberOfDatagramsIsNotPositive)
	})
}
