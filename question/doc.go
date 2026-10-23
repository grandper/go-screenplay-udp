// Package question holds the questions about the datagrams the actor
// received and about its connection:
//
//	datagrams, err := anActor.AsksFor(question.UDPDatagrams())                     // a []*ability.Datagram
//	datagrams, err := anActor.AsksFor(question.UDPDatagramsFrom("127.0.0.1:9000")) // a []*ability.Datagram
//	address, err := anActor.AsksFor(question.LocalUDPAddress())                    // a *net.UDPAddr
//
// A question composes with the core questions, the questions about a
// datagram and the resolutions of the extension:
//
//	err := anActor.Should(see.The(question.UDPDatagrams()).ShouldHave(been.Received()))
//	err := anActor.Should(see.The(text.Of(last.Of(question.UDPDatagrams()))).Should(be.EqualTo("pong")))
package question
