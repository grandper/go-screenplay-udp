// Package question gathers the questions of the extension. They live in
// packages named so that a statement reads as English:
//
//	datagrams, err := anActor.AsksFor(udp.Datagrams())                     // a []*udp.Datagram
//	datagrams, err := anActor.AsksFor(udp.DatagramsFrom("127.0.0.1:9000")) // a []*udp.Datagram
//	address, err := anActor.AsksFor(local.UDPAddress())                    // a *net.UDPAddr
//	text, err := anActor.AsksFor(text.Of(last.Of(udp.Datagrams())))        // a string
//
// A question about a datagram that is asked about a nil datagram fails with an
// error wrapping udp.ErrDatagramIsNil.
package question
