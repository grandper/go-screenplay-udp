# go-screenplay-udp

[![Test](https://github.com/grandper/go-screenplay-udp/actions/workflows/go-test.yml/badge.svg)](https://github.com/grandper/go-screenplay-udp/actions/workflows/go-test.yml)
[![Lint](https://github.com/grandper/go-screenplay-udp/actions/workflows/go-lint.yml/badge.svg)](https://github.com/grandper/go-screenplay-udp/actions/workflows/go-lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/grandper/go-screenplay-udp.svg)](https://pkg.go.dev/github.com/grandper/go-screenplay-udp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

This project is an extension of [go-screenplay](https://github.com/grandper/go-screenplay) to drive a UDP connection. It bundles the ability to communicate over UDP together with the actions, questions and resolutions that go with it, so you test a datagram scenario with the same screenplay vocabulary you use everywhere else. A test reads as English, and everything the actor does is narrated:
```go
actor := screenplay.ActorNamed("UDP User").WhoCan(communicate.OverUDP().ListeningOnAFreePort())

err := actor.AttemptsTo(send.TheDatagram("ping").To("127.0.0.1:9000"))
// UDP User attempts to send the datagram "ping" to "127.0.0.1:9000"
err = actor.AttemptsTo(receive.TheNextDatagram())
// UDP User attempts to receive the next datagram
err = actor.Should(see.The(udp.Datagrams()).ShouldHave(been.ReceivedFrom("127.0.0.1:9000")))
// UDP User should see that the UDP datagrams should have been received from "127.0.0.1:9000"
```

## Installation
```sh
go get github.com/grandper/go-screenplay-udp
```
It requires go-screenplay v2. The extension uses the `net` package of the standard library only.

## Ability
The `communicate` package provides the ability to communicate over UDP, which allows an actor to send and receive UDP datagrams through one `*net.UDPConn`. It is given to an actor with `communicate.OverUDP()`, reading "who can communicate over UDP". The ability is given an existing connection with `UsingTheConnection(conn)`, reading "communicate over UDP using the connection conn", an address to listen on with `ListeningOn("127.0.0.1:9000")`, reading "communicate over UDP listening on 127.0.0.1:9000", or a port chosen by the system on the loopback interface with `ListeningOnAFreePort()`, reading "communicate over UDP listening on a free port".

First, import the package:
```go
import (
    "github.com/grandper/go-screenplay/v2/screenplay"
    "github.com/grandper/go-screenplay/v2/action/see"
    "github.com/grandper/go-screenplay/v2/question/last"
    "github.com/grandper/go-screenplay/v2/resolution/be"
    "github.com/grandper/go-screenplay/v2/resolution/contain"
    "github.com/grandper/go-screenplay-udp/ability/communicate"
    "github.com/grandper/go-screenplay-udp/action/connect"
    "github.com/grandper/go-screenplay-udp/action/listen"
    "github.com/grandper/go-screenplay-udp/action/receive"
    "github.com/grandper/go-screenplay-udp/action/send"
    "github.com/grandper/go-screenplay-udp/action/stop"
    "github.com/grandper/go-screenplay-udp/fixture"
    "github.com/grandper/go-screenplay-udp/question/local"
    "github.com/grandper/go-screenplay-udp/question/payload"
    "github.com/grandper/go-screenplay-udp/question/sender"
    "github.com/grandper/go-screenplay-udp/question/text"
    "github.com/grandper/go-screenplay-udp/resolution/been"
    "github.com/grandper/go-screenplay-udp/udp"
)
```
Then create an actor who can communicate over UDP. The address is bound the first time the actor uses the ability, so an address that cannot be bound is reported by the first action, with an error wrapping `udp.ErrFailedToListen`:
```go
actor := screenplay.ActorNamed("UDP User").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
actor := screenplay.ActorNamed("UDP User").WhoCan(communicate.OverUDP().ListeningOn("127.0.0.1:9000"))
actor := screenplay.ActorNamed("UDP User").WhoCan(communicate.OverUDP().ListeningOn(":9000")) // every interface
```
Alternatively, give the ability a connection you opened yourself, with `net.ListenUDP`, `net.DialUDP` or `net.ListenMulticastUDP`:
```go
conn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9000})
actor := screenplay.ActorNamed("UDP User").WhoCan(communicate.OverUDP().UsingTheConnection(conn))
```
An actor without a connection or an address fails every action with `udp.ErrUDPConnectionNotSet`, "no UDP connection was given", until it performs one of the `listen` actions or `connect.ToTheUDPAddress`.

> **_NOTE:_** When the actor exits (`actor.Exit()`), the ability forgets what it holds: the received datagrams are dropped and the connection is closed, whether it was given with `UsingTheConnection` or bound by the ability itself.

`UsingTheConnection`, `ListeningOn` and `ListeningOnAFreePort` are the only methods a test can call on the ability. The `udp` package holds the `Datagram` the actor receives, next to the questions that answer with it, and every error the extension reports. Its operations live in `internal/ability`, so only an action or a question of this package uses them, in one statement written in `PerformAs` or `AnsweredBy` that reads "the ability of the actor to send":
```go
err := impl.TheAbilityOf(theActor).ToSend(data, "127.0.0.1:9000")
```
`TheAbilityOf` never fails: when the actor cannot communicate over UDP, the operation reports `screenplay.ErrActorMissingAbility`.

The ability is safe for concurrent use: actions running in parallel (see `proceed.ConcurrentlyTo`) can send and receive without racing. One datagram is received at a time, so concurrent receives wait for one another.

## Actions
The actions live in packages named after their verb, so that a statement reads as English: `send.TheDatagram("ping").To("127.0.0.1:9000")` reads "send the datagram ping to 127.0.0.1:9000". Every action is narrated word for word as it was written, after the name of the actor:
```go
err := actor.AttemptsTo(send.TheDatagram("ping").To("127.0.0.1:9000"))
// UDP User attempts to send the datagram "ping" to "127.0.0.1:9000"
```
Nothing panics when an action is built. A mistake, such as an empty address, is reported when the actor performs the action, with an error of the `udp` package: `udp.ErrAddressIsEmpty` or `udp.ErrNumberOfDatagramsIsNotPositive`. An address that cannot be resolved is reported with an error wrapping `udp.ErrInvalidAddress`.

### listen.OnTheUDPAddress
Binds the connection of the actor to an address, when its ability was given neither a connection nor an address. A port of 0 lets the system choose a free one:
```go
err := actor.AttemptsTo(listen.OnTheUDPAddress("127.0.0.1:9000"))
// UDP User attempts to listen on the UDP address "127.0.0.1:9000"
err := actor.AttemptsTo(listen.OnTheUDPAddress("127.0.0.1:0"))
```
`listen.OnAFreeUDPPort` binds a port chosen by the system on the loopback interface, which the question `local.UDPAddress` tells:
```go
err := actor.AttemptsTo(listen.OnAFreeUDPPort())
// UDP User attempts to listen on a free UDP port
```
They fail with `udp.ErrAddressIsEmpty` when the address is empty, with an error wrapping `udp.ErrFailedToListen` when the address cannot be bound, and with `udp.ErrAlreadyListening` when the actor already holds a connection.

### listen.OnTheMulticastGroup
Binds the connection of the actor to the address of a multicast group, on the multicast interface of the system, so that the actor receives the datagrams sent to the group:
```go
err := actor.AttemptsTo(listen.OnTheMulticastGroup("224.0.0.251:5353"))
// UDP User attempts to listen on the multicast group "224.0.0.251:5353"
```
It fails like `listen.OnTheUDPAddress`; an address that is not a multicast group cannot be joined, with an error wrapping `udp.ErrFailedToListen`.

### connect.ToTheUDPAddress
Connects the actor to a remote address: a datagram sent without an address goes there, and only the datagrams coming from there are received. A connected connection also reports a port nobody listens on, when the system does:
```go
err := actor.AttemptsTo(connect.ToTheUDPAddress("127.0.0.1:9000"))
// UDP User attempts to connect to the UDP address "127.0.0.1:9000"
err := actor.AttemptsTo(send.TheDatagram("ping"))
// UDP User attempts to send the datagram "ping"
```
It fails with `udp.ErrAddressIsEmpty` when the address is empty, with an error wrapping `udp.ErrFailedToConnect` when the connection cannot be opened, and with `udp.ErrAlreadyListening` when the actor already holds a connection.

### send.TheDatagram
Sends a datagram to the address given with `To`. The datagram holds a text (`TheDatagram`), bytes (`TheBinaryDatagram`) or nothing (`AnEmptyDatagram`):
```go
err := actor.AttemptsTo(send.TheDatagram("ping").To("127.0.0.1:9000"))
// UDP User attempts to send the datagram "ping" to "127.0.0.1:9000"
err := actor.AttemptsTo(send.TheBinaryDatagram([]byte{0x00, 0x01}).To("127.0.0.1:9000"))
// UDP User attempts to send the binary datagram "\x00\x01" to "127.0.0.1:9000"
err := actor.AttemptsTo(send.AnEmptyDatagram().To("127.0.0.1:9000"))
// UDP User attempts to send an empty datagram to "127.0.0.1:9000"
```
`ToTheSenderOfTheLastDatagram` replies to whoever sent the last datagram the actor received, which is how an actor plays a UDP server:
```go
err := actor.AttemptsTo(send.TheDatagram("pong").ToTheSenderOfTheLastDatagram())
// UDP User attempts to send the datagram "pong" to the sender of the last datagram
```
A broadcast is a datagram sent to the broadcast address of the network:
```go
err := actor.AttemptsTo(send.TheDatagram("anyone?").To("255.255.255.255:9000"))
```
It fails with `udp.ErrAddressIsEmpty` when `To` is given an empty address, with `udp.ErrDestinationNotSet` when no address is given and the actor is not connected to one, with an error wrapping `udp.ErrConnectedToAnotherAddress` when the address differs from the one the actor connected to, with `udp.ErrNoDatagramReceived` when replying while nothing was received, and with an error wrapping `udp.ErrFailedToSend` when the system refuses the datagram, for example one larger than 65,507 bytes.

### receive.TheNextDatagram
Receives the next datagram and stores it, where `udp.Datagrams()` finds it. The actor waits at most its timeout (two seconds by default, see `screenplay.WithTimeout`), unless `TimingOutAfter` is written:
```go
err := actor.AttemptsTo(receive.TheNextDatagram())
// UDP User attempts to receive the next datagram
err := actor.AttemptsTo(receive.TheNextDatagram().TimingOutAfter(500).Milliseconds())
// UDP User attempts to receive the next datagram timing out after 500 milliseconds
```
`receive.TheNext(amount).Datagrams()` receives several datagrams, waiting at most the timeout for all of them:
```go
err := actor.AttemptsTo(receive.TheNext(3).Datagrams())
// UDP User attempts to receive the next 3 datagrams
err := actor.AttemptsTo(receive.TheNext(1).Datagram().TimingOutAfter(2).Seconds())
// UDP User attempts to receive the next 1 datagram timing out after 2 seconds
```
A datagram that does not come in time fails the action with an error wrapping `udp.ErrFailedToReceive` and `context.DeadlineExceeded`; the datagrams received before it are kept. A cancellation of the actor's context stops the wait as well, with `context.Canceled`. It fails with `udp.ErrNumberOfDatagramsIsNotPositive` when the number of datagrams is not positive.

Datagrams that arrive while the actor is not receiving wait in the buffer of the socket, so they are received by a later action. A datagram larger than the buffer of the socket is truncated by the system, as UDP does.

### stop.ListeningOverUDP
Closes the connection of the actor and forgets the address it was bound to, so that the actor can listen again elsewhere. The datagrams received so far are kept:
```go
err := actor.AttemptsTo(stop.ListeningOverUDP())
// UDP User attempts to stop listening over UDP
```
It fails with `udp.ErrUDPConnectionNotSet` when the actor holds no connection.

## Questions
A question is a noun phrase that the actor asks for, and its answer is typed: `see.The`, `AsksFor` and the core questions such as `last.Of` and `number.Of` work on it without any type assertion.

### udp.Datagrams
`udp.Datagrams()` answers with the `[]*udp.Datagram` received by the actor, in the order they were received. `udp.DatagramsFrom(address)` keeps the datagrams sent from that address, resolved when the question is answered, so `"localhost:9000"` matches the datagrams sent from `127.0.0.1:9000`:
```go
err := actor.Should(see.The(udp.Datagrams()).ShouldHave(been.Received()))
// UDP User should see that the UDP datagrams should have been received
err := actor.Should(see.The(udp.DatagramsFrom("127.0.0.1:9000")).ShouldHave(been.Received()))
// UDP User should see that the UDP datagrams from "127.0.0.1:9000" should have been received
err := actor.Should(see.The(number.Of(udp.Datagrams())).Should(be.EqualTo(2)))
lastDatagram, err := actor.AsksFor(last.Of(udp.Datagrams())) // lastDatagram is an *udp.Datagram
```
A datagram is read through `Data()`, a copy of its payload, `Text()`, the payload as a text, and `Sender()`, a copy of the `*net.UDPAddr` it was sent from.

### The parts of a datagram
`text.Of`, `payload.Of` and `sender.Of` take a question that answers with one datagram, such as `last.Of(udp.Datagrams())`, and answer with a part of it, so that a core resolution inspects it:
```go
err := actor.Should(see.The(text.Of(last.Of(udp.Datagrams()))).Should(be.EqualTo("pong")))
// UDP User should see that the text of last of UDP datagrams should be equal to "pong"
err := actor.Should(see.The(text.Of(last.Of(udp.Datagrams()))).Should(contain.Substring("pong")))
data, err := actor.AsksFor(payload.Of(last.Of(udp.Datagrams())))    // data is a []byte
address, err := actor.AsksFor(sender.Of(last.Of(udp.Datagrams()))) // address is a *net.UDPAddr
```
A nil datagram is an error wrapping `udp.ErrDatagramIsNil`.

### local.UDPAddress
`local.UDPAddress()` answers with a copy of the `*net.UDPAddr` the connection of the actor is bound to, which tells the port the system chose when the actor listens on a free port:
```go
address, err := actor.AsksFor(local.UDPAddress()) // address is a *net.UDPAddr
// UDP User asks for the local UDP address
```

### Fake answers
The `question/fixture` package builds the answer of `udp.Datagrams()` without a UDP peer, to test code that receives it. An address that cannot be resolved leaves the sender unknown:
```go
pong := fixture.CreateFakeDatagram("pong", "127.0.0.1:9000")
raw := fixture.CreateFakeBinaryDatagram([]byte{0x00, 0x01}, "127.0.0.1:9000")
```

## Resolutions
The `been` package holds the resolutions of this extension. Each one reads after the verb `ShouldHave`, and a negation lives only in the verb: `ShouldNotHave(been.Received())` reads "should not have been received".

`been.Received()` and `been.ReceivedFrom(address)` inspect the answer of `udp.Datagrams()` or `udp.DatagramsFrom(address)`. The address of `ReceivedFrom` is resolved when the datagrams are inspected, so `"localhost:9000"` matches a datagram sent from `127.0.0.1:9000`; an address that cannot be resolved makes the resolution invalid, with an error wrapping `screenplay.ErrInvalidResolution`:
```go
err := actor.Should(see.The(udp.Datagrams()).ShouldHave(been.Received()))
// UDP User should see that the UDP datagrams should have been received
// fails: UDP User failed to see that the UDP datagrams should have been received: no datagram has been received
err := actor.Should(see.The(udp.Datagrams()).ShouldHave(been.ReceivedFrom("127.0.0.1:9000")))
// UDP User should see that the UDP datagrams should have been received from "127.0.0.1:9000"
err := actor.Should(see.The(udp.DatagramsFrom("127.0.0.1:9000")).ShouldNotHave(been.Received()))
// UDP User should see that the UDP datagrams from "127.0.0.1:9000" should not have been received
```

## Testing with a UDP peer
The `fixture` package runs in-process UDP peers so scenarios can run without an external system. Every helper takes a `*testing.T` and registers its own cleanup, so the peer is stopped when the test ends:
```go
// Runs a server on a free port of the loopback interface that sends every
// datagram back to its sender, and returns its address.
echo := fixture.RunAnEchoUDPServer(t)

// Runs a server that replies with what the function returns, or does not
// reply when it returns nil.
server := fixture.RunAUDPServer(t, func(data []byte) []byte { return []byte("pong") })

// Binds a connection to a free port of the loopback interface, to give to
// UsingTheConnection or to play a peer by hand.
conn := fixture.ListenOnAFreeUDPPort(t)
```

## Putting it all together
A complete scenario, from running the echo server to asserting on the reply:
```go
func TestReceivesTheEchoOfADatagram(t *testing.T) {
    echo := fixture.RunAnEchoUDPServer(t)
    actor := screenplay.ActorNamed("UDP User").WhoCan(communicate.OverUDP().ListeningOnAFreePort())

    require.NoError(t, actor.AttemptsTo(
        send.TheDatagram("ping").To(echo.String()),
        receive.TheNextDatagram(),
    ))

    require.NoError(t, actor.Should(
        see.The(udp.Datagrams()).ShouldHave(been.ReceivedFrom(echo.String())),
        see.The(text.Of(last.Of(udp.Datagrams()))).Should(be.EqualTo("ping")),
    ))
}
```
And two actors talking to each other, one playing the server:
```go
func TestAServerRepliesToAClient(t *testing.T) {
    server := screenplay.ActorNamed("Server").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
    client := screenplay.ActorNamed("Client").WhoCan(communicate.OverUDP().ListeningOnAFreePort())
    serverAddress, err := server.AsksFor(local.UDPAddress())
    require.NoError(t, err)

    require.NoError(t, client.AttemptsTo(send.TheDatagram("ping").To(serverAddress.String())))
    require.NoError(t, server.AttemptsTo(
        receive.TheNextDatagram(),
        send.TheDatagram("pong").ToTheSenderOfTheLastDatagram(),
    ))
    require.NoError(t, client.AttemptsTo(receive.TheNextDatagram()))

    require.NoError(t, client.Should(
        see.The(text.Of(last.Of(udp.DatagramsFrom(serverAddress.String())))).Should(be.EqualTo("pong")),
    ))
}
```

## References
- [Package net](https://pkg.go.dev/net#UDPConn) of the Go standard library
- [RFC 768](https://www.rfc-editor.org/rfc/rfc768), the User Datagram Protocol
