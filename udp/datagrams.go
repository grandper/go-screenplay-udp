package udp

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	impl "github.com/grandper/go-screenplay-udp/internal/ability"
	"github.com/grandper/go-screenplay-udp/internal/verdict"
)

// Datagrams asks for every datagram the actor received, in the order they
// were received.
func Datagrams() screenplay.Question[[]*Datagram] {
	return &datagramsQuestion{}
}

// DatagramsFrom asks for the datagrams the actor received from the
// address, in the order they were received. The address is resolved when the
// question is answered, so "localhost:9000" matches the datagrams sent from
// 127.0.0.1:9000.
func DatagramsFrom(address string) screenplay.Question[[]*Datagram] {
	return &datagramsQuestion{address: address, fromAddress: true}
}

// datagramsQuestion asks for the datagrams the actor received.
type datagramsQuestion struct {
	address     string
	fromAddress bool
}

// AnsweredBy returns the datagrams the actor received.
func (q *datagramsQuestion) AnsweredBy(theActor *screenplay.Actor) ([]*Datagram, error) {
	if q.fromAddress {
		return impl.TheAbilityOf(theActor).ToRetrieveDatagramsFrom(q.address)
	}

	return impl.TheAbilityOf(theActor).ToRetrieveDatagrams()
}

// String describes the question.
func (q *datagramsQuestion) String() string {
	if q.fromAddress {
		return "UDP datagrams from " + verdict.Describe(q.address)
	}

	return "UDP datagrams"
}

// datagramsQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]*Datagram] = (*datagramsQuestion)(nil)
