package ability

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/grandper/go-screenplay/v2/screenplay"
)

// maxDatagramSize is the largest payload a UDP datagram can carry: the
// receive buffer is that large, so no datagram is truncated.
const maxDatagramSize = 65535

// FreePortAddress is the address bound by ListeningOnAFreePort and by the
// action listen.OnAFreeUDPPort: a port chosen by the system on the loopback
// interface.
const FreePortAddress = "127.0.0.1:0"

var (
	// ErrUDPConnectionNotSet is returned when the ability was given neither a
	// UDP connection nor an address to listen on.
	ErrUDPConnectionNotSet = errors.New("no UDP connection was given")

	// ErrFailedToListen is wrapped by the error returned when the UDP address
	// cannot be bound.
	ErrFailedToListen = errors.New("failed to listen on the UDP address")

	// ErrFailedToConnect is wrapped by the error returned when the UDP
	// connection to the remote address cannot be opened.
	ErrFailedToConnect = errors.New("failed to connect to the UDP address")

	// ErrAlreadyListening is returned when the actor listens or connects while
	// its ability already holds a connection.
	ErrAlreadyListening = errors.New("the actor already holds a UDP connection")

	// ErrInvalidAddress is wrapped by the error returned when an address cannot
	// be resolved to a UDP address.
	ErrInvalidAddress = errors.New("invalid UDP address")

	// ErrDestinationNotSet is returned when a datagram is sent without an
	// address, over a connection that is not connected to a remote address.
	ErrDestinationNotSet = errors.New("the datagram has no destination")

	// ErrConnectedToAnotherAddress is wrapped by the error returned when a
	// datagram is sent to an address over a connection connected to another
	// one.
	ErrConnectedToAnotherAddress = errors.New("the connection is connected to another address")

	// ErrFailedToSend is wrapped by the error returned when a datagram cannot
	// be sent.
	ErrFailedToSend = errors.New("failed to send the datagram")

	// ErrFailedToReceive is wrapped by the error returned when a datagram
	// cannot be received, including when none arrives in time.
	ErrFailedToReceive = errors.New("failed to receive a datagram")

	// ErrNoDatagramReceived is returned when the actor replies to the sender of
	// the last datagram while it has received none.
	ErrNoDatagramReceived = errors.New("no datagram has been received")
)

// CommunicateOverUDP enables the actor to communicate over UDP. The
// connection is given with UsingTheConnection, or bound by the ability itself
// to the address given with ListeningOn or ListeningOnAFreePort, reading
// "communicate over UDP using the connection conn" and "communicate over UDP
// listening on 127.0.0.1:9000":
//
//	communicate.OverUDP().UsingTheConnection(conn)
//	communicate.OverUDP().ListeningOn("127.0.0.1:9000")
//	communicate.OverUDP().ListeningOnAFreePort()
func CommunicateOverUDP() *CommunicateOverUDPAbility {
	return &CommunicateOverUDPAbility{}
}

// CommunicateOverUDPAbility is the ability to communicate over UDP. It is safe
// for concurrent use: actions running in parallel (see proceed.ConcurrentlyTo)
// can send and receive without racing, one datagram being received at a time.
// Its method set holds only its settings and Forget; the operations are
// reached through TheAbilityOf, so a test cannot bypass the actor.
type CommunicateOverUDPAbility struct {
	// binding serializes the lazy binding so that the state mutex is never
	// held during a network call.
	binding sync.Mutex
	address string

	// receiving serializes the reads, as a read deadline is set on the
	// connection for each of them.
	receiving sync.Mutex

	mutex    sync.RWMutex
	conn     *net.UDPConn
	received []*Datagram
}

// UsingTheConnection provides the UDP connection to the ability; it reads
// "using the connection conn". A connection already held by the ability is
// kept, and a nil connection is ignored.
func (c *CommunicateOverUDPAbility) UsingTheConnection(conn *net.UDPConn) *CommunicateOverUDPAbility {
	if c == nil || conn == nil {
		return c
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn == nil {
		c.conn = conn
	}

	return c
}

// ListeningOn gives the ability the address to listen on; it reads
// "listening on 127.0.0.1:9000". The address is bound the first time the
// ability is used, so that a failure to bind is reported by the action that
// needed it. A connection already held by the ability is kept.
func (c *CommunicateOverUDPAbility) ListeningOn(address string) *CommunicateOverUDPAbility {
	if c == nil {
		return c
	}

	c.binding.Lock()
	defer c.binding.Unlock()

	c.address = address

	return c
}

// ListeningOnAFreePort gives the ability a port chosen by the system on the
// loopback interface to listen on; it reads "listening on a free port". The
// port is bound the first time the ability is used, and the question
// local.UDPAddress tells which one it is.
func (c *CommunicateOverUDPAbility) ListeningOnAFreePort() *CommunicateOverUDPAbility {
	return c.ListeningOn(FreePortAddress)
}

// Forget forgets the received datagrams and closes the connection. On a nil
// receiver, the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (c *CommunicateOverUDPAbility) Forget() error {
	if c == nil {
		return screenplay.ErrActorMissingAbility
	}

	c.binding.Lock()
	defer c.binding.Unlock()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.address = ""
	c.received = nil
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil

	return err
}

// String describes the ability.
func (c *CommunicateOverUDPAbility) String() string {
	return "communicate over UDP"
}

// CommunicateOverUDPAbility implements the screenplay.Ability interface.
var _ screenplay.Ability = (*CommunicateOverUDPAbility)(nil)

// TheAbilityOf returns the operations of the actor's ability to communicate
// over UDP; it is written in PerformAs or AnsweredBy and reads "the ability
// of the actor to send":
//
//	err := ability.TheAbilityOf(theActor).ToSend(data, "127.0.0.1:9000")
//
// It never fails. When the actor is nil or does not have the ability, every
// operation fails with screenplay.ErrActorMissingAbility.
func TheAbilityOf(actor *screenplay.Actor) *Operations {
	return &Operations{ability: screenplay.UseAnAbilityOf(actor).That(CommunicateOverUDP())}
}

// Operations holds the operations of the actor's ability to communicate over
// UDP. They are not methods of CommunicateOverUDPAbility, whose method set
// holds only its settings, so a test cannot send a datagram without an action.
type Operations struct {
	ability *CommunicateOverUDPAbility
}

// ToListenOn binds the connection of the ability to the address. It fails
// with ErrAlreadyListening when the ability already holds a connection, with
// an error wrapping ErrInvalidAddress when the address cannot be resolved,
// and with an error wrapping ErrFailedToListen when it cannot be bound. When
// the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToListenOn(address string) error {
	return o.open(func(c *CommunicateOverUDPAbility) error { return c.listen(address) })
}

// ToListenOnTheMulticastGroup binds the connection of the ability to the
// address of a multicast group, on the multicast interface of the system, so
// that it receives the datagrams sent to the group. It fails like ToListenOn.
// When the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToListenOnTheMulticastGroup(address string) error {
	return o.open(func(c *CommunicateOverUDPAbility) error { return c.listenMulticast(address) })
}

// ToConnectTo opens a connection to the remote address, over which a datagram
// sent without an address goes to that address and only the datagrams it
// sends are received. It fails with ErrAlreadyListening when the ability
// already holds a connection, with an error wrapping ErrInvalidAddress when
// the address cannot be resolved, and with an error wrapping ErrFailedToConnect
// when the connection cannot be opened. When the actor does not have the
// ability, it fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToConnectTo(address string) error {
	return o.open(func(c *CommunicateOverUDPAbility) error { return c.connect(address) })
}

// open runs the function that opens the connection, unless the ability
// already holds one.
func (o *Operations) open(opening func(*CommunicateOverUDPAbility) error) error {
	c := o.ability
	if c == nil {
		return screenplay.ErrActorMissingAbility
	}

	c.binding.Lock()
	defer c.binding.Unlock()

	if c.isOpen() {
		return ErrAlreadyListening
	}

	return opening(c)
}

// ToSend sends the datagram to the address, or to the remote address of a
// connected connection when the address is empty. It fails with
// ErrDestinationNotSet when there is neither, with an error wrapping
// ErrConnectedToAnotherAddress when the address differs from the remote
// address of a connected connection, and with an error wrapping
// ErrFailedToSend when the datagram cannot be sent. When the actor does not
// have the ability, it fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToSend(data []byte, address string) error {
	conn, err := o.connection()
	if err != nil {
		return err
	}

	if address == "" {
		if conn.RemoteAddr() == nil {
			return ErrDestinationNotSet
		}

		return send(data, func() (int, error) { return conn.Write(data) })
	}

	destination, err := resolve(address)
	if err != nil {
		return err
	}
	if remote, ok := conn.RemoteAddr().(*net.UDPAddr); ok && remote != nil {
		if remote.String() != destination.String() {
			return fmt.Errorf(
				"%w: cannot send to %q over the connection to %q",
				ErrConnectedToAnotherAddress,
				address,
				remote,
			)
		}

		return send(data, func() (int, error) { return conn.Write(data) })
	}

	return send(data, func() (int, error) { return conn.WriteToUDP(data, destination) })
}

// ToSendToTheSenderOfTheLastDatagram sends the datagram to the address the
// last received datagram came from. It fails with ErrNoDatagramReceived when
// the actor has received none. When the actor does not have the ability, it
// fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToSendToTheSenderOfTheLastDatagram(data []byte) error {
	if _, err := o.connection(); err != nil {
		return err
	}
	c := o.ability

	c.mutex.RLock()
	last := len(c.received) - 1
	var sender *net.UDPAddr
	if last >= 0 {
		sender = c.received[last].Sender()
	}
	c.mutex.RUnlock()

	if sender == nil {
		return ErrNoDatagramReceived
	}

	return o.ToSend(data, sender.String())
}

// ToReceiveTheNextDatagram waits for the next datagram, until the context
// ends, and stores it. A datagram that does not come in time is reported with
// an error wrapping ErrFailedToReceive and the error of the context. When the
// actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToReceiveTheNextDatagram(ctx context.Context) error {
	conn, err := o.connection()
	if err != nil {
		return err
	}
	c := o.ability

	c.receiving.Lock()
	defer c.receiving.Unlock()

	datagram, err := receive(ctx, conn)
	if err != nil {
		return err
	}
	c.store(datagram)

	return nil
}

// ToRetrieveDatagrams returns a copy of every datagram received, in the order
// they were received, whether or not the ability holds a connection: the
// datagrams are kept after the connection is closed. When the actor does not
// have the ability, it fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToRetrieveDatagrams() ([]*Datagram, error) {
	return o.retrieveDatagrams(func(*Datagram) bool { return true })
}

// ToRetrieveDatagramsFrom returns a copy of the datagrams received from the
// address, in the order they were received. It fails with an error wrapping
// ErrInvalidAddress when the address cannot be resolved. When the actor does
// not have the ability, it fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToRetrieveDatagramsFrom(address string) ([]*Datagram, error) {
	sender, err := resolve(address)
	if err != nil {
		return nil, err
	}

	return o.retrieveDatagrams(func(datagram *Datagram) bool { return datagram.IsFrom(sender) })
}

// retrieveDatagrams returns a copy of the wanted datagrams. It does not need
// the connection, so that it neither binds the address as a side effect nor
// fails once the connection is closed.
func (o *Operations) retrieveDatagrams(wanted func(*Datagram) bool) ([]*Datagram, error) {
	c := o.ability
	if c == nil {
		return nil, screenplay.ErrActorMissingAbility
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	datagrams := make([]*Datagram, 0, len(c.received))
	for _, datagram := range c.received {
		if wanted(datagram) {
			datagrams = append(datagrams, datagram)
		}
	}

	return datagrams, nil
}

// ToTellTheLocalAddress returns a copy of the address the connection is bound
// to. When the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToTellTheLocalAddress() (*net.UDPAddr, error) {
	conn, err := o.connection()
	if err != nil {
		return nil, err
	}

	local, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok || local == nil {
		return nil, fmt.Errorf("%w: the connection has no local address", ErrUDPConnectionNotSet)
	}

	return cloneAddress(local), nil
}

// ToClose closes the connection and forgets the address it was bound to, so
// that the actor can listen again. The received datagrams are kept. When the
// actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToClose() error {
	c := o.ability
	if c == nil {
		return screenplay.ErrActorMissingAbility
	}

	c.binding.Lock()
	defer c.binding.Unlock()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.address = ""
	if c.conn == nil {
		return ErrUDPConnectionNotSet
	}
	err := c.conn.Close()
	c.conn = nil

	return err
}

// connection returns the connection of the ability, binding it to the address
// given with ListeningOn when needed. It fails with
// screenplay.ErrActorMissingAbility when the actor does not have the ability,
// and with ErrUDPConnectionNotSet when the ability has neither a connection
// nor an address.
func (o *Operations) connection() (*net.UDPConn, error) {
	c := o.ability
	if c == nil {
		return nil, screenplay.ErrActorMissingAbility
	}

	c.binding.Lock()
	defer c.binding.Unlock()

	if !c.isOpen() {
		if c.address == "" {
			return nil, ErrUDPConnectionNotSet
		}
		if err := c.listen(c.address); err != nil {
			return nil, err
		}
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.conn, nil
}

// listen binds a connection to the address and stores it. The binding mutex
// must be held.
func (c *CommunicateOverUDPAbility) listen(address string) error {
	local, err := resolve(address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", local)
	if err != nil {
		return fmt.Errorf("%w %q: %w", ErrFailedToListen, address, err)
	}
	c.use(conn)

	return nil
}

// listenMulticast binds a connection to the multicast group at the address
// and stores it. The binding mutex must be held.
func (c *CommunicateOverUDPAbility) listenMulticast(address string) error {
	group, err := resolve(address)
	if err != nil {
		return err
	}

	conn, err := net.ListenMulticastUDP("udp", nil, group)
	if err != nil {
		return fmt.Errorf("%w %q: %w", ErrFailedToListen, address, err)
	}
	c.use(conn)

	return nil
}

// connect opens a connection to the remote address and stores it. The
// binding mutex must be held.
func (c *CommunicateOverUDPAbility) connect(address string) error {
	remote, err := resolve(address)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, remote)
	if err != nil {
		return fmt.Errorf("%w %q: %w", ErrFailedToConnect, address, err)
	}
	c.use(conn)

	return nil
}

// use stores the connection.
func (c *CommunicateOverUDPAbility) use(conn *net.UDPConn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.conn = conn
}

func (c *CommunicateOverUDPAbility) isOpen() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.conn != nil
}

// store records a received datagram.
func (c *CommunicateOverUDPAbility) store(datagram *Datagram) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.received = append(c.received, datagram)
}

// resolve turns the address into a UDP address, failing with an error wrapping
// ErrInvalidAddress when it cannot.
func resolve(address string) (*net.UDPAddr, error) {
	resolved, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("%w %q: %w", ErrInvalidAddress, address, err)
	}

	return resolved, nil
}

// send writes the datagram with the function and reports a failure or a
// partial write with an error wrapping ErrFailedToSend.
func send(data []byte, write func() (int, error)) error {
	written, err := write()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSend, err)
	}
	if written != len(data) {
		return fmt.Errorf("%w: %d of %d bytes were written", ErrFailedToSend, written, len(data))
	}

	return nil
}

// receive reads the next datagram of the connection, until the context ends.
// The deadline of the context bounds the read, and its cancellation ends it.
func receive(ctx context.Context, conn *net.UDPConn) (*Datagram, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToReceive, err)
	}

	deadline, _ := ctx.Deadline()
	if err := conn.SetReadDeadline(deadline); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToReceive, err)
	}
	stop := context.AfterFunc(ctx, func() { _ = conn.SetReadDeadline(time.Now()) })
	defer stop()

	buffer := make([]byte, maxDatagramSize)
	read, sender, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, receiveError(ctx, err)
	}

	return NewDatagram(buffer[:read], sender), nil
}

// receiveError explains why no datagram was received: the context ended, the
// deadline passed, or the read failed.
func receiveError(ctx context.Context, err error) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return fmt.Errorf("%w: %w: %w", ErrFailedToReceive, ctxErr, err)
	}
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return fmt.Errorf("%w: %w: %w", ErrFailedToReceive, context.DeadlineExceeded, err)
	}

	return fmt.Errorf("%w: %w", ErrFailedToReceive, err)
}
