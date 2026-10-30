// Package udp holds the questions about the datagrams the actor received,
// which read "the UDP datagrams" and "the UDP datagrams from 127.0.0.1:9000",
// the Datagram they answer with, and every error the extension reports:
//
//	datagrams, err := anActor.AsksFor(udp.Datagrams())                     // a []*udp.Datagram
//	datagrams, err := anActor.AsksFor(udp.DatagramsFrom("127.0.0.1:9000")) // a []*udp.Datagram
//	errors.Is(err, udp.ErrFailedToListen)
//
// A datagram is read through Data(), Text() and Sender(). The datagrams kept
// are bounded by the buffered size of the actor's production. A question composes
// with the core questions, the questions about a datagram and the resolutions
// of the extension:
//
//	err := anActor.Should(see.The(udp.Datagrams()).ShouldHave(been.Received()))
//	err := anActor.Should(see.The(text.Of(last.Of(udp.Datagrams()))).Should(be.EqualTo("pong")))
//
// The ability itself is given to an actor with communicate.OverUDP.
package udp
