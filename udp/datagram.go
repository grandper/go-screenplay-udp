package udp

import impl "github.com/grandper/go-screenplay-udp/internal/ability"

// Datagram is a datagram received by the actor, which Datagrams and
// DatagramsFrom answer with: its payload and the address it was sent from,
// read through Data(), Text() and Sender(). It does not change once received.
type Datagram = impl.Datagram
