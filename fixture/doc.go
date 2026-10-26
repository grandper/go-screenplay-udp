// Package fixture holds the UDP peers that let a scenario run without an
// external system: a server that echoes every datagram it receives, and a
// connection bound to a free port.
//
//	echo := fixture.RunAnEchoUDPServer(t)
//	anActor := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
//	err := anActor.AttemptsTo(send.TheDatagram("ping").To(echo.String()))
package fixture
