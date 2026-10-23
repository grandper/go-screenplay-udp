package question

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/ability"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
	"github.com/grandper/go-screenplay-udp/internal/verdict"
)

// UDPDatagrams asks for every datagram the actor received, in the order they
// were received.
func UDPDatagrams() screenplay.Question[[]*ability.Datagram] {
	return &udpDatagramsQuestion{}
}

// UDPDatagramsFrom asks for the datagrams the actor received from the
// address, in the order they were received. The address is resolved when the
// question is answered, so "localhost:9000" matches the datagrams sent from
// 127.0.0.1:9000.
func UDPDatagramsFrom(address string) screenplay.Question[[]*ability.Datagram] {
	return &udpDatagramsQuestion{address: address, fromAddress: true}
}

// udpDatagramsQuestion asks for the datagrams the actor received.
type udpDatagramsQuestion struct {
	address     string
	fromAddress bool
}

// AnsweredBy returns the datagrams the actor received.
func (q *udpDatagramsQuestion) AnsweredBy(theActor *screenplay.Actor) ([]*ability.Datagram, error) {
	if q.fromAddress {
		return impl.TheAbilityOf(theActor).ToRetrieveDatagramsFrom(q.address)
	}

	return impl.TheAbilityOf(theActor).ToRetrieveDatagrams()
}

// String describes the question.
func (q *udpDatagramsQuestion) String() string {
	if q.fromAddress {
		return "UDP datagrams from " + verdict.Describe(q.address)
	}

	return "UDP datagrams"
}

// udpDatagramsQuestion implements the screenplay.Question interface.
var _ screenplay.Question[[]*ability.Datagram] = (*udpDatagramsQuestion)(nil)
