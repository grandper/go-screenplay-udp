// Package communicate holds the ability to communicate over UDP, which reads
// "who can communicate over UDP listening on a free port":
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
//
// The ability is given an existing connection with UsingTheConnection, or an
// address to listen on with ListeningOn, reading "communicate over UDP
// listening on 127.0.0.1:9000":
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().UsingTheConnection(conn))
//	anActor := screenplay.ActorNamed("Adam").WhoCan(communicate.OverUDP().ListeningOn("127.0.0.1:9000"))
package communicate
