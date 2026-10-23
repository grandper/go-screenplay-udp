package ability

import impl "github.com/grandper/go-screenplay-udp/internal/ability"

// CommunicateOverUDP enables the actor to communicate over UDP. The
// connection is given with UsingTheConnection, or bound by the ability itself
// to the address given with ListeningOn or ListeningOnAFreePort, reading
// "communicate over UDP using the connection conn" and "communicate over UDP
// listening on 127.0.0.1:9000":
//
//	CommunicateOverUDP().UsingTheConnection(conn)
//	CommunicateOverUDP().ListeningOn("127.0.0.1:9000")
//	CommunicateOverUDP().ListeningOnAFreePort()
func CommunicateOverUDP() *impl.CommunicateOverUDPAbility {
	return impl.CommunicateOverUDP()
}
