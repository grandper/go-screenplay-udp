package send

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/action"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
	"github.com/grandper/go-screenplay-udp/internal/verdict"
)

// TheDatagram sends a datagram holding the text to the address given with To,
// or to the sender of the last received datagram with
// ToTheSenderOfTheLastDatagram; it reads "send the datagram ping to
// 127.0.0.1:9000". Without either, the datagram goes to the address the actor
// connected to:
//
//	send.TheDatagram("ping").To("127.0.0.1:9000")
//	send.TheDatagram("pong").ToTheSenderOfTheLastDatagram()
//	send.TheDatagram("ping")
func TheDatagram(text string) *TheDatagramAction {
	return &TheDatagramAction{
		description: "send the datagram " + verdict.Describe(text),
		data:        []byte(text),
	}
}

// TheBinaryDatagram sends a datagram holding the bytes to the address given
// with To; it reads "send the binary datagram data to 127.0.0.1:9000":
//
//	send.TheBinaryDatagram(data).To("127.0.0.1:9000")
func TheBinaryDatagram(data []byte) *TheDatagramAction {
	return &TheDatagramAction{
		description: "send the binary datagram " + verdict.Describe(data),
		data:        data,
	}
}

// AnEmptyDatagram sends a datagram with no payload to the address given with
// To; it reads "send an empty datagram to 127.0.0.1:9000":
//
//	send.AnEmptyDatagram().To("127.0.0.1:9000")
func AnEmptyDatagram() *TheDatagramAction {
	return &TheDatagramAction{description: "send an empty datagram"}
}

// TheDatagramAction sends a datagram.
type TheDatagramAction struct {
	description string
	data        []byte
	address     string
	toWritten   bool
	toTheSender bool
}

// To sets the address the datagram is sent to; it reads "to 127.0.0.1:9000".
func (a *TheDatagramAction) To(address string) *TheDatagramAction {
	a.address = address
	a.toWritten = true
	a.toTheSender = false

	return a
}

// ToTheSenderOfTheLastDatagram sends the datagram to the address the last
// received datagram came from; it reads "to the sender of the last datagram".
// It fails with ability.ErrNoDatagramReceived when the actor has received no
// datagram.
func (a *TheDatagramAction) ToTheSenderOfTheLastDatagram() *TheDatagramAction {
	a.address = ""
	a.toWritten = false
	a.toTheSender = true

	return a
}

// PerformAs performs the action as the provided actor. It fails with
// action.ErrAddressIsEmpty when To was given an empty address, and with
// ability.ErrDestinationNotSet when no address was given and the actor is not
// connected to one.
func (a *TheDatagramAction) PerformAs(theActor *screenplay.Actor) error {
	if a.toTheSender {
		return impl.TheAbilityOf(theActor).ToSendToTheSenderOfTheLastDatagram(a.data)
	}
	if a.toWritten && a.address == "" {
		return action.ErrAddressIsEmpty
	}

	return impl.TheAbilityOf(theActor).ToSend(a.data, a.address)
}

// String describes the action, word for word as it was written.
func (a *TheDatagramAction) String() string {
	switch {
	case a.toTheSender:
		return a.description + " to the sender of the last datagram"
	case a.toWritten:
		return a.description + " to " + verdict.Describe(a.address)
	default:
		return a.description
	}
}

// TheDatagramAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheDatagramAction)(nil)
