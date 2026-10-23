package question

import "errors"

// ErrDatagramIsNil is wrapped by the error returned when a question about a
// datagram is asked about a nil datagram.
var ErrDatagramIsNil = errors.New("the datagram is nil")
