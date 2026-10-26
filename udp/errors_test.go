package udp_test

import (
	"testing"

	"github.com/grandper/go-screenplay/v2/screenplay"
	screenplayfixture "github.com/grandper/go-screenplay/v2/screenplay/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability/communicate"
	"github.com/grandper/go-screenplay-udp/action/listen"
	"github.com/grandper/go-screenplay-udp/action/receive"
	"github.com/grandper/go-screenplay-udp/action/send"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
	"github.com/grandper/go-screenplay-udp/question/payload"
	"github.com/grandper/go-screenplay-udp/question/sender"
	"github.com/grandper/go-screenplay-udp/question/text"
	"github.com/grandper/go-screenplay-udp/udp"
)

func TestErrors(t *testing.T) {
	ofTheAbility := map[error]error{
		udp.ErrUDPConnectionNotSet:       impl.ErrUDPConnectionNotSet,
		udp.ErrFailedToListen:            impl.ErrFailedToListen,
		udp.ErrFailedToConnect:           impl.ErrFailedToConnect,
		udp.ErrAlreadyListening:          impl.ErrAlreadyListening,
		udp.ErrInvalidAddress:            impl.ErrInvalidAddress,
		udp.ErrDestinationNotSet:         impl.ErrDestinationNotSet,
		udp.ErrConnectedToAnotherAddress: impl.ErrConnectedToAnotherAddress,
		udp.ErrFailedToSend:              impl.ErrFailedToSend,
		udp.ErrFailedToReceive:           impl.ErrFailedToReceive,
		udp.ErrNoDatagramReceived:        impl.ErrNoDatagramReceived,
	}
	all := []error{
		udp.ErrAddressIsEmpty,
		udp.ErrNumberOfDatagramsIsNotPositive,
		udp.ErrDatagramIsNil,
	}
	for public := range ofTheAbility {
		all = append(all, public)
	}

	t.Run("of the ability are the errors of the internal ability", func(t *testing.T) {
		for public, internal := range ofTheAbility {
			require.ErrorIs(t, internal, public)
			assert.Same(t, internal, public)
		}
	})

	t.Run("are distinct from one another", func(t *testing.T) {
		for _, err := range all {
			for _, other := range all {
				if err.Error() != other.Error() {
					require.NotErrorIs(t, err, other)
				}
			}
		}
	})

	t.Run("describe the failure", func(t *testing.T) {
		require.EqualError(t, udp.ErrUDPConnectionNotSet, "no UDP connection was given")
		require.EqualError(t, udp.ErrFailedToListen, "failed to listen on the UDP address")
		require.EqualError(t, udp.ErrFailedToConnect, "failed to connect to the UDP address")
		require.EqualError(t, udp.ErrAlreadyListening, "the actor already holds a UDP connection")
		require.EqualError(t, udp.ErrInvalidAddress, "invalid UDP address")
		require.EqualError(t, udp.ErrDestinationNotSet, "the datagram has no destination")
		require.EqualError(t, udp.ErrConnectedToAnotherAddress, "the connection is connected to another address")
		require.EqualError(t, udp.ErrFailedToSend, "failed to send the datagram")
		require.EqualError(t, udp.ErrFailedToReceive, "failed to receive a datagram")
		require.EqualError(t, udp.ErrNoDatagramReceived, "no datagram has been received")
		require.EqualError(t, udp.ErrAddressIsEmpty, "the UDP address must be provided")
		require.EqualError(
			t,
			udp.ErrNumberOfDatagramsIsNotPositive,
			"the number of datagrams to receive must be positive",
		)
		require.EqualError(t, udp.ErrDatagramIsNil, "the datagram is nil")
	})

	t.Run("of the actions are reported when the actor performs the action, not when it is built", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP())

		var sendToNowhere, listenNowhere, receiveNothing screenplay.Performable
		require.NotPanics(t, func() {
			sendToNowhere = send.TheDatagram("ping").To("")
			listenNowhere = listen.OnTheUDPAddress("")
			receiveNothing = receive.TheNext(0).Datagrams()
		})

		require.ErrorIs(t, adam.AttemptsTo(sendToNowhere), udp.ErrAddressIsEmpty)
		require.ErrorIs(t, adam.AttemptsTo(listenNowhere), udp.ErrAddressIsEmpty)
		require.ErrorIs(t, adam.AttemptsTo(receiveNothing), udp.ErrNumberOfDatagramsIsNotPositive)
	})

	t.Run("ErrDatagramIsNil is wrapped by the error of a question asked about a nil datagram", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")
		noDatagram := screenplayfixture.NewFakeQuestion[*udp.Datagram]("the datagram", nil)

		_, err := adam.AsksFor(payload.Of(noDatagram))
		require.ErrorIs(t, err, udp.ErrDatagramIsNil)
		_, err = adam.AsksFor(text.Of(noDatagram))
		require.ErrorIs(t, err, udp.ErrDatagramIsNil)
		_, err = adam.AsksFor(sender.Of(noDatagram))
		require.ErrorIs(t, err, udp.ErrDatagramIsNil)
	})
}
