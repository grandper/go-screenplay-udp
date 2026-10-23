package action

import "errors"

var (
	// ErrAddressIsEmpty is returned when an action is given an empty address.
	// Nothing panics when the action is built; the mistake is reported when the
	// actor performs it.
	ErrAddressIsEmpty = errors.New("the UDP address must be provided")

	// ErrNumberOfDatagramsIsNotPositive is returned when receive.TheNext is
	// given a number of datagrams that is not positive. Nothing panics when the
	// action is built; the mistake is reported when the actor performs it.
	ErrNumberOfDatagramsIsNotPositive = errors.New("the number of datagrams to receive must be positive")
)
