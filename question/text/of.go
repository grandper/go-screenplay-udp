package text

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/internal/question/datagram"
	"github.com/grandper/go-screenplay-udp/udp"
)

// Of asks for the payload of the datagram the question answers with, as a
// text; it reads "text of last of UDP datagrams". A nil datagram is an error
// wrapping udp.ErrDatagramIsNil.
func Of(question screenplay.Question[*udp.Datagram]) screenplay.Question[string] {
	return datagram.Of("text of", question, (*udp.Datagram).Text)
}
