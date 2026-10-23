// Package stop holds the action that closes the UDP connection of the actor,
// so that it can listen again elsewhere:
//
//	err := anActor.AttemptsTo(stop.ListeningOverUDP())
package stop
