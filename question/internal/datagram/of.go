package datagram

import (
	"fmt"

	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/ability"
	"github.com/grandper/go-screenplay-udp/question"
)

// Of creates a question that answers with the part of the datagram the
// function reads, after the wrapped question answered with the datagram; the
// description is the head of the phrase, such as "text of". A nil datagram
// is an error wrapping question.ErrDatagramIsNil.
func Of[T any](
	head string,
	wrapped screenplay.Question[*ability.Datagram],
	read func(*ability.Datagram) T,
) screenplay.Question[T] {
	return &ofQuestion[T]{head: head, wrapped: wrapped, read: read}
}

// ofQuestion asks for a part of the datagram the wrapped question answers
// with.
type ofQuestion[T any] struct {
	head    string
	wrapped screenplay.Question[*ability.Datagram]
	read    func(*ability.Datagram) T
}

// AnsweredBy asks the wrapped question and reads the part of its answer.
func (q *ofQuestion[T]) AnsweredBy(theActor *screenplay.Actor) (T, error) {
	var zero T

	datagram, err := q.wrapped.AnsweredBy(theActor)
	if err != nil {
		return zero, err
	}
	if datagram == nil {
		return zero, fmt.Errorf("%w: %s", question.ErrDatagramIsNil, q.wrapped)
	}

	return q.read(datagram), nil
}

// String describes the question.
func (q *ofQuestion[T]) String() string {
	return fmt.Sprintf("%s %s", q.head, q.wrapped)
}
