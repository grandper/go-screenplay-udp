package ability

import impl "github.com/grandper/go-screenplay-udp/internal/ability"

var (
	// ErrUDPConnectionNotSet is returned when the ability was given neither a
	// UDP connection nor an address to listen on.
	ErrUDPConnectionNotSet = impl.ErrUDPConnectionNotSet

	// ErrFailedToListen is wrapped by the error returned when the UDP address
	// cannot be bound.
	ErrFailedToListen = impl.ErrFailedToListen

	// ErrFailedToConnect is wrapped by the error returned when the UDP
	// connection to the remote address cannot be opened.
	ErrFailedToConnect = impl.ErrFailedToConnect

	// ErrAlreadyListening is returned when the actor listens or connects while
	// its ability already holds a connection.
	ErrAlreadyListening = impl.ErrAlreadyListening

	// ErrInvalidAddress is wrapped by the error returned when an address cannot
	// be resolved to a UDP address.
	ErrInvalidAddress = impl.ErrInvalidAddress

	// ErrDestinationNotSet is returned when a datagram is sent without an
	// address, over a connection that is not connected to a remote address.
	ErrDestinationNotSet = impl.ErrDestinationNotSet

	// ErrConnectedToAnotherAddress is wrapped by the error returned when a
	// datagram is sent to an address over a connection connected to another
	// one.
	ErrConnectedToAnotherAddress = impl.ErrConnectedToAnotherAddress

	// ErrFailedToSend is wrapped by the error returned when a datagram cannot
	// be sent.
	ErrFailedToSend = impl.ErrFailedToSend

	// ErrFailedToReceive is wrapped by the error returned when a datagram
	// cannot be received, including when none arrives in time.
	ErrFailedToReceive = impl.ErrFailedToReceive

	// ErrNoDatagramReceived is returned when the actor replies to the sender of
	// the last datagram while it has received none.
	ErrNoDatagramReceived = impl.ErrNoDatagramReceived
)
