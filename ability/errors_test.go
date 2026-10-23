package ability_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay-udp/ability"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
)

func TestErrors(t *testing.T) {
	errs := map[error]error{
		ability.ErrUDPConnectionNotSet:       impl.ErrUDPConnectionNotSet,
		ability.ErrFailedToListen:            impl.ErrFailedToListen,
		ability.ErrFailedToConnect:           impl.ErrFailedToConnect,
		ability.ErrAlreadyListening:          impl.ErrAlreadyListening,
		ability.ErrInvalidAddress:            impl.ErrInvalidAddress,
		ability.ErrDestinationNotSet:         impl.ErrDestinationNotSet,
		ability.ErrConnectedToAnotherAddress: impl.ErrConnectedToAnotherAddress,
		ability.ErrFailedToSend:              impl.ErrFailedToSend,
		ability.ErrFailedToReceive:           impl.ErrFailedToReceive,
		ability.ErrNoDatagramReceived:        impl.ErrNoDatagramReceived,
	}

	t.Run("are the errors of the internal ability", func(t *testing.T) {
		for public, internal := range errs {
			require.ErrorIs(t, internal, public)
			assert.Same(t, internal, public)
		}
	})

	t.Run("are distinct from one another", func(t *testing.T) {
		for public := range errs {
			for other := range errs {
				if public.Error() != other.Error() {
					require.NotErrorIs(t, public, other)
				}
			}
		}
	})

	t.Run("describe the failure", func(t *testing.T) {
		require.EqualError(t, ability.ErrUDPConnectionNotSet, "no UDP connection was given")
		require.EqualError(t, ability.ErrFailedToListen, "failed to listen on the UDP address")
		require.EqualError(t, ability.ErrFailedToConnect, "failed to connect to the UDP address")
		require.EqualError(t, ability.ErrAlreadyListening, "the actor already holds a UDP connection")
		require.EqualError(t, ability.ErrInvalidAddress, "invalid UDP address")
		require.EqualError(t, ability.ErrDestinationNotSet, "the datagram has no destination")
		require.EqualError(t, ability.ErrConnectedToAnotherAddress, "the connection is connected to another address")
		require.EqualError(t, ability.ErrFailedToSend, "failed to send the datagram")
		require.EqualError(t, ability.ErrFailedToReceive, "failed to receive a datagram")
		require.EqualError(t, ability.ErrNoDatagramReceived, "no datagram has been received")
	})
}
