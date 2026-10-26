package communicate

import impl "github.com/grandper/go-screenplay-udp/internal/ability"

// OverUDP enables the actor to communicate over UDP. The connection is given
// with UsingTheConnection, or bound by the ability itself to the address
// given with ListeningOn or ListeningOnAFreePort, reading "communicate over
// UDP using the connection conn" and "communicate over UDP listening on
// 127.0.0.1:9000":
//
//	communicate.OverUDP().UsingTheConnection(conn)
//	communicate.OverUDP().ListeningOn("127.0.0.1:9000")
//	communicate.OverUDP().ListeningOnAFreePort()
func OverUDP() *impl.CommunicateOverUDPAbility {
	return impl.CommunicateOverUDP()
}
