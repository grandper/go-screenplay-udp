// Package listen holds the actions that bind the connection of the actor to
// a UDP address, when its ability was given neither a connection nor an
// address:
//
//	err := anActor.AttemptsTo(listen.OnTheUDPAddress("127.0.0.1:9000"))
//	err := anActor.AttemptsTo(listen.OnAFreeUDPPort())
//	err := anActor.AttemptsTo(listen.OnTheMulticastGroup("224.0.0.251:5353"))
package listen
