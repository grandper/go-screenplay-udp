package sender

import (
	"net"

	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/question/internal/datagram"
)

// Of asks for the address the datagram the question answers with was sent
// from; it reads "sender of last of UDP datagrams". A nil datagram is an
// error wrapping question.ErrDatagramIsNil.
func Of(question screenplay.Question[*ability.Datagram]) screenplay.Question[*net.UDPAddr] {
	return datagram.Of("sender of", question, (*ability.Datagram).Sender)
}
