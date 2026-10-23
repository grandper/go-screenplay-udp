// Package ability holds the ability to communicate over UDP:
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOnAFreePort())
//
// The ability is given an existing connection with UsingTheConnection, or an
// address to listen on with ListeningOn, reading "communicate over UDP
// listening on 127.0.0.1:9000":
//
//	anActor := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().UsingTheConnection(conn))
//	anActor := screenplay.ActorNamed("Adam").WhoCan(ability.CommunicateOverUDP().ListeningOn("127.0.0.1:9000"))
package ability
