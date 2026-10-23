package receive

import (
	"context"
	"fmt"
	"time"

	"github.com/grandper/go-screenplay/v2/action/timing"
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/action"
	impl "github.com/grandper/go-screenplay-udp/internal/ability"
)

// TheNextDatagram receives the next datagram and stores it, where the question
// UDPDatagrams finds it. The actor waits at most its timeout (see
// screenplay.WithTimeout), unless TimingOutAfter is written; it reads
// "receive the next datagram timing out after 2 seconds":
//
//	receive.TheNextDatagram()
//	receive.TheNextDatagram().TimingOutAfter(2).Seconds()
func TheNextDatagram() *TheNextDatagramAction {
	return &TheNextDatagramAction{amount: 1, description: "receive the next datagram"}
}

// TheNext receives the next datagrams, as many as the amount, and stores them;
// it reads "receive the next 3 datagrams". The actor waits at most its timeout
// for all of them, unless TimingOutAfter is written:
//
//	receive.TheNext(3).Datagrams()
//	receive.TheNext(3).Datagrams().TimingOutAfter(2).Seconds()
func TheNext(amount int) *TheNextBuilder {
	return &TheNextBuilder{amount: amount}
}

// TheNextBuilder holds the amount given to TheNext and waits for the noun
// that closes the statement.
type TheNextBuilder struct {
	amount int
}

// Datagrams makes the amount the number of datagrams to receive.
func (b *TheNextBuilder) Datagrams() *TheNextDatagramAction {
	unit := "datagrams"
	if b.amount == 1 {
		unit = "datagram"
	}

	return &TheNextDatagramAction{
		amount:      b.amount,
		description: fmt.Sprintf("receive the next %d %s", b.amount, unit),
	}
}

// Datagram is the singular of Datagrams, so that one datagram reads as
// English.
func (b *TheNextBuilder) Datagram() *TheNextDatagramAction {
	return b.Datagrams()
}

// TheNextDatagramAction receives the next datagrams.
type TheNextDatagramAction struct {
	amount         int
	description    string
	timeout        time.Duration
	timeoutBuilder *timing.DurationBuilder[TheNextDatagramAction]
	timeoutWritten bool
}

// TimingOutAfter bounds the wait for the datagrams; it reads "timing out
// after 2 seconds". Without it, the actor waits at most its timeout. A
// datagram that does not come in time fails the action with an error wrapping
// ability.ErrFailedToReceive and context.DeadlineExceeded.
func (a *TheNextDatagramAction) TimingOutAfter(amount int) *timing.DurationBuilder[TheNextDatagramAction] {
	a.timeoutWritten = true
	a.timeoutBuilder = timing.NewDurationBuilder(a, &a.timeout).For(amount)

	return a.timeoutBuilder
}

// PerformAs performs the action as the provided actor. It fails with
// action.ErrNumberOfDatagramsIsNotPositive when the number of datagrams is
// not positive.
func (a *TheNextDatagramAction) PerformAs(theActor *screenplay.Actor) error {
	if a.amount <= 0 {
		return action.ErrNumberOfDatagramsIsNotPositive
	}

	timeout := theActor.Timeout()
	if a.timeoutWritten {
		timeout = a.timeout
	}
	ctx, cancel := context.WithTimeout(theActor.Context(), timeout)
	defer cancel()

	for range a.amount {
		if err := impl.TheAbilityOf(theActor).ToReceiveTheNextDatagram(ctx); err != nil {
			return err
		}
	}

	return nil
}

// String describes the action, word for word as it was written.
func (a *TheNextDatagramAction) String() string {
	if a.timeoutWritten {
		return a.description + " timing out after " + a.timeoutBuilder.String()
	}

	return a.description
}

// TheNextDatagramAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheNextDatagramAction)(nil)
