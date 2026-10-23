// Package receive holds the actions that receive the next datagrams and store
// them, where the question UDPDatagrams finds them:
//
//	err := anActor.AttemptsTo(receive.TheNextDatagram())
//	err := anActor.AttemptsTo(receive.TheNextDatagram().TimingOutAfter(500).Milliseconds())
//	err := anActor.AttemptsTo(receive.TheNext(3).Datagrams())
package receive
