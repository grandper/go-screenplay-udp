package sender

import (
	"net"

	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/internal/question/datagram"
	"github.com/grandper/go-screenplay-udp/udp"
)

// Of asks for the address the datagram the question answers with was sent
// from; it reads "sender of last of UDP datagrams". A nil datagram is an
// error wrapping udp.ErrDatagramIsNil.
func Of(question screenplay.Question[*udp.Datagram]) screenplay.Question[*net.UDPAddr] {
	return datagram.Of("sender of", question, (*udp.Datagram).Sender)
}
