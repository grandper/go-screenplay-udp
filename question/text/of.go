package text

import (
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/question/internal/datagram"
)

// Of asks for the payload of the datagram the question answers with, as a
// text; it reads "text of last of UDP datagrams". A nil datagram is an error
// wrapping question.ErrDatagramIsNil.
func Of(question screenplay.Question[*ability.Datagram]) screenplay.Question[string] {
	return datagram.Of("text of", question, (*ability.Datagram).Text)
}
