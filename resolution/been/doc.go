// Package been holds the resolutions about the datagrams the actor received.
// Each one reads after the verb ShouldHave:
//
//	err := anActor.Should(see.The(udp.Datagrams()).ShouldHave(been.Received()))
//	err := anActor.Should(see.The(udp.Datagrams()).ShouldHave(been.ReceivedFrom("127.0.0.1:9000")))
//
// They read "see that the UDP datagrams should have been received" and "see
// that the UDP datagrams should have been received from 127.0.0.1:9000".
package been
