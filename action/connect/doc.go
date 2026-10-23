// Package connect holds the action that connects the actor to a remote UDP
// address, so that a datagram sent without an address goes there and only
// the datagrams coming from there are received:
//
//	err := anActor.AttemptsTo(connect.ToTheUDPAddress("127.0.0.1:9000"))
package connect
