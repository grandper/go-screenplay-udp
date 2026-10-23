// Package send holds the actions that send a datagram:
//
//	err := anActor.AttemptsTo(send.TheDatagram("ping").To("127.0.0.1:9000"))
//	err := anActor.AttemptsTo(send.TheBinaryDatagram(data).To("127.0.0.1:9000"))
//	err := anActor.AttemptsTo(send.AnEmptyDatagram().To("127.0.0.1:9000"))
//	err := anActor.AttemptsTo(send.TheDatagram("pong").ToTheSenderOfTheLastDatagram())
//
// Over a connection opened with connect.ToTheUDPAddress, a datagram is sent
// without an address:
//
//	err := anActor.AttemptsTo(send.TheDatagram("ping"))
package send
