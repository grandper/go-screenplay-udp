package receive

import impl "github.com/grandper/go-screenplay-udp/action/internal/receive"

// TheNextDatagram receives the next datagram and stores it, where the question
// UDPDatagrams finds it. The actor waits at most its timeout (see
// screenplay.WithTimeout), unless TimingOutAfter is written; it reads
// "receive the next datagram timing out after 2 seconds":
//
//	receive.TheNextDatagram()
//	receive.TheNextDatagram().TimingOutAfter(2).Seconds()
func TheNextDatagram() *impl.TheNextDatagramAction {
	return impl.TheNextDatagram()
}

// TheNext receives the next datagrams, as many as the amount, and stores them;
// it reads "receive the next 3 datagrams". The actor waits at most its timeout
// for all of them, unless TimingOutAfter is written:
//
//	receive.TheNext(3).Datagrams()
//	receive.TheNext(3).Datagrams().TimingOutAfter(2).Seconds()
func TheNext(amount int) *impl.TheNextBuilder {
	return impl.TheNext(amount)
}
