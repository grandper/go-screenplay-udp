// Package fixture holds the fakes to test code that receives the answer of
// the UDP datagrams questions, without a UDP peer:
//
//	pong := fixture.CreateFakeDatagram("pong", "127.0.0.1:9000")
//	fixture.AssertMatch(t, been.ReceivedFrom("127.0.0.1:9000"), []*ability.Datagram{pong})
package fixture
