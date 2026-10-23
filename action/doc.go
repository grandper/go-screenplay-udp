// Package action holds the errors of the actions that drive a UDP
// connection. The actions live in packages named after their verb:
//
//	listen.OnTheUDPAddress("127.0.0.1:9000")              action/listen
//	listen.OnAFreeUDPPort()                                action/listen
//	listen.OnTheMulticastGroup("224.0.0.251:5353")         action/listen
//	connect.ToTheUDPAddress("127.0.0.1:9000")              action/connect
//	send.TheDatagram("ping").To("127.0.0.1:9000")          action/send
//	receive.TheNextDatagram()                              action/receive
//	receive.TheNext(3).Datagrams()                         action/receive
//	stop.ListeningOverUDP()                                action/stop
//
// For example:
//
//	err := anActor.AttemptsTo(
//		listen.OnAFreeUDPPort(),
//		send.TheDatagram("ping").To("127.0.0.1:9000"),
//		receive.TheNextDatagram(),
//	)
package action
