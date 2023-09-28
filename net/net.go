package net

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Addr represents a network end point address.
//
// The two methods Network and String conventionally return strings
// that can be passed as the arguments to Dial, but the exact form
// and meaning of the strings is up to the implementation.
type Addr interface {
	Network() string // name of the network (for example, "tcp", "udp")
	String() string  // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
}

type Conn interface {
	// Write writes data to the connection.
	// Write can be made to time out and return an error after a fixed
	// time limit; see SetDeadline and SetWriteDeadline.
	Write(b []byte) (n int, err error)
}

type conn struct {
	fd *netFD
}

func (c *conn) ok() bool { return c != nil && c.fd != nil }

var listenerBacklogCache struct {
	sync.Once
	val int
}

// listenerBacklog is a caching wrapper around maxListenerBacklog.
func listenerBacklog() int {
	listenerBacklogCache.Do(func() { listenerBacklogCache.val = maxListenerBacklog() })
	return listenerBacklogCache.val
}

// Various errors contained in OpError.
var (
	// For connection setup and write operations.
	errMissingAddress = errors.New("missing address")

	// For both read and write operations.
	errCanceled = errors.New("operation was canceled")
)

// mapErr maps from the context errors to the historical internal net
// error values.
//
// TODO(bradfitz): get rid of this after adjusting tests and making
// context.DeadlineExceeded implement net.Error?
func mapErr(err error) error {
	switch err {
	case context.Canceled:
		return errCanceled
	case context.DeadlineExceeded:
		return errTimeout
	default:
		return err
	}
}

// OpError is the error type usually returned by functions in the net
// package. It describes the operation, network type, and address of
// an error.
type OpError struct {
	// Op is the operation which caused the error, such as
	// "read" or "write".
	Op string

	// Net is the network type on which this error occurred,
	// such as "tcp" or "udp6".
	Net string

	// For operations involving a remote network connection, like
	// Dial, Read, or Write, Source is the corresponding local
	// network address.
	Source Addr

	// Addr is the network address for which this error occurred.
	// For local operations, like Listen or SetDeadline, Addr is
	// the address of the local endpoint being manipulated.
	// For operations involving a remote network connection, like
	// Dial, Read, or Write, Addr is the remote address of that
	// connection.
	Addr Addr

	// Err is the error that occurred during the operation.
	// The Error method panics if the error is nil.
	Err error
}

func (e *OpError) Unwrap() error { return e.Err }

func (e *OpError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := e.Op
	if e.Net != "" {
		s += " " + e.Net
	}
	if e.Source != nil {
		s += " " + e.Source.String()
	}
	if e.Addr != nil {
		if e.Source != nil {
			s += "->"
		} else {
			s += " "
		}
		s += e.Addr.String()
	}
	s += ": " + e.Err.Error()
	return s
}

var (
	// aLongTimeAgo is a non-zero time, far in the past, used for
	// immediate cancellation of dials.
	aLongTimeAgo = time.Unix(1, 0)

	// nonDeadline and noCancel are just zero values for
	// readability with functions taking too many parameters.
	noDeadline = time.Time{}
	noCancel   = (chan struct{})(nil)
)

type AddrError struct {
	Err  string
	Addr string
}

func (e *AddrError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := e.Err
	if e.Addr != "" {
		s = "address " + e.Addr + ": " + s
	}
	return s
}

func (e *AddrError) Timeout() bool   { return false }
func (e *AddrError) Temporary() bool { return false }

// errTimeout exists to return the historical "i/o timeout" string
// for context.DeadlineExceeded. See mapErr.
// It is also used when Dialer.Deadline is exceeded.
//
// TODO(iant): We could consider changing this to os.ErrDeadlineExceeded
// in the future, but note that that would conflict with the TODO
// at mapErr that suggests changing it to context.DeadlineExceeded.
var errTimeout error = &timeoutError{}

type timeoutError struct{}

func (e *timeoutError) Error() string   { return "i/o timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }
