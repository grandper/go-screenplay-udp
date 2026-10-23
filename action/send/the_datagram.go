package send

import impl "github.com/grandper/go-screenplay-udp/action/internal/send"

// TheDatagram sends a datagram holding the text to the address given with To,
// or to the sender of the last received datagram with
// ToTheSenderOfTheLastDatagram; it reads "send the datagram ping to
// 127.0.0.1:9000". Without either, the datagram goes to the address the actor
// connected to:
//
//	send.TheDatagram("ping").To("127.0.0.1:9000")
//	send.TheDatagram("pong").ToTheSenderOfTheLastDatagram()
//	send.TheDatagram("ping")
func TheDatagram(text string) *impl.TheDatagramAction {
	return impl.TheDatagram(text)
}

// TheBinaryDatagram sends a datagram holding the bytes to the address given
// with To; it reads "send the binary datagram data to 127.0.0.1:9000":
//
//	send.TheBinaryDatagram(data).To("127.0.0.1:9000")
func TheBinaryDatagram(data []byte) *impl.TheDatagramAction {
	return impl.TheBinaryDatagram(data)
}

// AnEmptyDatagram sends a datagram with no payload to the address given with
// To; it reads "send an empty datagram to 127.0.0.1:9000":
//
//	send.AnEmptyDatagram().To("127.0.0.1:9000")
func AnEmptyDatagram() *impl.TheDatagramAction {
	return impl.AnEmptyDatagram()
}
