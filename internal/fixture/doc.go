// Package fixture holds the helpers of the extension's own tests: an actor
// whose narration is recorded, fake datagrams to build the answer of
// udp.Datagrams without a UDP peer, and in-process UDP peers so that a
// scenario runs without an external system:
//
//	anActor := fixture.NarratedActor("Adam", recorder, fixture.ListenOnAFreeUDPPort(t))
//	pong := fixture.CreateFakeDatagram("pong", "127.0.0.1:9000")
//	echo := fixture.RunAnEchoUDPServer(t)
package fixture
