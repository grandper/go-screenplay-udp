package been

import (
	"errors"
	"fmt"
	"net"

	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/screenplay"

	"github.com/grandper/go-screenplay-udp/internal/verdict"
	"github.com/grandper/go-screenplay-udp/udp"
)

// Received creates a resolution to tell if at least one datagram has been
// received; it reads "been received".
func Received() screenplay.Resolution[[]*udp.Datagram] {
	return resolution.FromFunc("been received", func(datagrams []*udp.Datagram) error {
		if len(datagrams) == 0 {
			return errors.New("no datagram has been received")
		}

		return nil
	})
}

// ReceivedFrom creates a resolution to tell if at least one datagram has been
// received from the address; it reads "been received from 127.0.0.1:9000".
// The address is resolved when the resolution inspects the datagrams, so
// "localhost:9000" matches a datagram sent from 127.0.0.1:9000; an address
// that cannot be resolved makes the resolution invalid, with an error wrapping
// screenplay.ErrInvalidResolution.
func ReceivedFrom(address string) screenplay.Resolution[[]*udp.Datagram] {
	return resolution.FromFunc(
		"been received from "+verdict.Describe(address),
		func(datagrams []*udp.Datagram) error {
			sender, err := net.ResolveUDPAddr("udp", address)
			if err != nil {
				return fmt.Errorf("%w: %w", screenplay.ErrInvalidResolution, err)
			}
			for _, datagram := range datagrams {
				if datagram != nil && datagram.IsFrom(sender) {
					return nil
				}
			}

			return fmt.Errorf("no datagram has been received from %s", verdict.Describe(address))
		},
	)
}
